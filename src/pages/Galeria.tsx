import { useState, useEffect } from "react";
import { Calendar, MapPin, Users, Image as ImageIcon, ImagePlus, Plus, X, Trash2, AlertCircle, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ShowWhenAuthenticated } from "@/auth/AuthSwitch";
import { FileUploadModal } from "@/components/FileUploadModal";
import { useToast } from "@/hooks/use-toast";
import * as api from "@/services/api";

const GallerySlug = "gallery_mais_fotos"

interface Photo {
  id: string;
  url: string;
  caption: string;
  uploadedAt: string;
}

interface GaleryEventDisplay {
  id: string;
  name: string;
  date: string;
  location: string;
  photos: Photo[];
}

type SelectedPhotoSource =
  | { source: "gallery" }
  | { source: "event"; eventId: string };

type SelectedPhoto = Photo &
  SelectedPhotoSource & {
    eventName?: string;
    eventDate?: string;
    eventLocation?: string;
  };

const Galeria = () => {
  const { toast } = useToast();
  const [events, setEvents] = useState<GaleryEventDisplay[]>([]);
  const [allPhotos, setAllPhotos] = useState<(Photo & { eventName?: string; eventDate?: string; eventLocation?: string })[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [selectedPhoto, setSelectedPhoto] = useState<SelectedPhoto | null>(null);
  const [showImageUpload, setShowImageUpload] = useState(false);
  const [showCreateEvent, setShowCreateEvent] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [eventToDelete, setEventToDelete] = useState<GaleryEventDisplay | null>(null);
  const [isDeletingEvent, setIsDeletingEvent] = useState(false);
  const [isCreatingEvent, setIsCreatingEvent] = useState(false);
  const [newEventData, setNewEventData] = useState({
    name: "",
    date: "",
    location: "",
    description: "",
  });
  const [newEventImages, setNewEventImages] = useState<File[]>([]);

  // Fetch galery events and images on mount
  useEffect(() => {
    fetchGalleryData();
  }, []);

  const fetchGalleryData = async () => {
    try {
      setLoading(true);
      setError(null);

      // Fetch galery events for "Eventos" section
      const galeryEvents = await api.listGaleryEvents();

      // Transform galery events to display format
      const transformedEvents: GaleryEventDisplay[] = galeryEvents.map(event => ({
        id: event.id,
        name: event.name,
        date: event.date,
        location: event.location,
        photos: event.image_urls.map((url, index) => ({
          // Use real Firestore ID if available, otherwise fall back to composite ID for old events
          id: event.image_ids?.[index] || `${event.id}-${index}`,
          url,
          caption: `${event.name} - Foto ${index + 1}`,
          uploadedAt: event.created_at,
        })),
      }));

      setEvents(transformedEvents);

      // Also fetch all individual images from the backend
      try {
        const galleryImages = await api.getImagesBySlug(GallerySlug);

        // Transform gallery images to photo format
        const transformedGalleryImages = galleryImages.map(img => ({
          id: img.id,
          url: img.object_url,
          caption: img.name,
          uploadedAt: img.created_at,
          eventName: img.text || img.slug || 'Galeria Geral',
          eventDate: img.date || img.created_at,
          eventLocation: img.location,
        }));

        // Sort all photos by upload date (newest first)
        transformedGalleryImages.sort((a, b) =>
          new Date(b.uploadedAt).getTime() - new Date(a.uploadedAt).getTime()
        );

        setAllPhotos(transformedGalleryImages);
      } catch (imgErr) {
        // If gallery images don't exist yet, that's okay - just continue with event photos
        console.log('No gallery images found or error fetching them:', imgErr);
      }


    } catch (err) {
      console.error('Failed to fetch gallery data:', err);
      setError('Falha ao carregar as fotos. Tente novamente mais tarde.');
    } finally {
      setLoading(false);
    }
  };

  // Refresh all images (optimized for individual image uploads)
  const refreshGalleryImages = async () => {
    try {
      // Fetch all individual images from the backend
      const galleryImages = await api.getImagesBySlug(GallerySlug);

      // Transform gallery images to photo format
      const transformedGalleryImages = galleryImages.map(img => ({
        id: img.id,
        url: img.object_url,
        caption: img.name,
        uploadedAt: img.created_at,
        eventName: img.text || img.slug || 'Galeria Geral',
        eventDate: img.date || img.created_at,
        eventLocation: img.location,
      }));

      // Get existing event photos (keep them as is)
      const allEventPhotos = events.flatMap(event =>
        event.photos.map(photo => ({
          ...photo,
          eventName: event.name,
          eventDate: event.date,
          eventLocation: event.location,
        }))
      );

      // Merge event photos with refreshed gallery images
      const mergedPhotos = [...allEventPhotos, ...transformedGalleryImages];

      // Sort by upload date (newest first)
      mergedPhotos.sort((a, b) =>
        new Date(b.uploadedAt).getTime() - new Date(a.uploadedAt).getTime()
      );

      setAllPhotos(mergedPhotos);
    } catch (err) {
      console.error('Failed to refresh gallery images:', err);
      // If refresh fails, fall back to full refresh
      await fetchGalleryData();
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('pt-BR');
  };

  const handleUploadImage = async (files: File[]) => {
    try {
      console.log('Uploading individual images:', files);

      // Upload each image to the /images endpoint
      const uploadPromises = files.map(async (file) => {
        // Convert image to base64
        const base64Data = await api.fileToBase64(file);

        // Create image with gallery slug "galeria-geral" (general gallery)
        return api.createImage({
          slug: GallerySlug,
          name: file.name,
          text: `Imagem adicionada em ${new Date().toLocaleDateString('pt-BR')}`,
          data: base64Data,
        });
      });

      // Wait for all uploads to complete
      const uploadedImages = await Promise.all(uploadPromises);
      console.log('Images uploaded successfully:', uploadedImages);

      // Only refresh the galeria-geral images, not all events
      await refreshGalleryImages();

      toast({
        title: "Sucesso!",
        description: uploadedImages.length === 1
          ? "1 imagem enviada com sucesso!"
          : `${uploadedImages.length} imagens enviadas com sucesso!`,
      });
    } catch (err) {
      console.error('Failed to upload images:', err);
      toast({
        variant: "destructive",
        title: "Erro ao enviar imagens",
        description: "Falha ao enviar imagens. Verifique o console para mais detalhes.",
      });
    }
  };

  const handleDeleteImage = async () => {
    if (!selectedPhoto) return;

    try {
      setIsDeleting(true);
      setDeleteError(null);

      await api.deleteImage(selectedPhoto.id);

      // Close modal
      setSelectedPhoto(null);

      // Refresh the gallery to remove the deleted image
      await fetchGalleryData();
    } catch (err) {
      console.error('Failed to delete image:', err);
      setDeleteError('Falha ao excluir a imagem. Tente novamente.');
    } finally {
      setIsDeleting(false);
    }
  };

  const handleDeleteEventPhoto = async (photo: SelectedPhoto) => {
    if (photo.source !== "event") return;

    try {
      setIsDeleting(true);
      setDeleteError(null);

      const event = events.find((e) => e.id === photo.eventId);
      if (!event) {
        throw new Error("Evento não encontrado");
      }

      const remainingPhotos = event.photos.filter((p) => p.id !== photo.id);
      console.log(remainingPhotos)

      const request: api.ModifyGaleryEventRequest = {
        id: event.id,
        name: event.name,
        location: event.location,
        date: event.date,
        image_ids: remainingPhotos.map((p) => p.id),
        image_urls: remainingPhotos.map((p) => p.url),
      };

      await api.modifyGaleryEventRequest(request);

      // Close modal and refresh
      setSelectedPhoto(null);
      await fetchGalleryData();
    } catch (err) {
      console.error("Failed to remove photo from event:", err);
      setDeleteError("Falha ao excluir a imagem do evento. Tente novamente.");
    } finally {
      setIsDeleting(false);
    }
  };


  // Check if a photo can be deleted (has a real Firestore ID, not a composite ID)
  const canDeletePhoto = (photoId: string): boolean => {
    // Composite IDs have format: {eventId}-{index}, e.g., "abc123-0"
    // Real Firestore IDs don't have this pattern (they're alphanumeric without dashes at the end)
    return !/.*-\d+$/.test(photoId);
  };

  const handleDeleteEvent = async () => {
    if (!eventToDelete) return;

    try {
      setIsDeletingEvent(true);
      setDeleteError(null);

      await api.deleteGaleryEvent(eventToDelete.id);

      // Close modal
      setEventToDelete(null);

      // Refresh the gallery to remove the deleted event
      await fetchGalleryData();
    } catch (err) {
      console.error('Failed to delete event:', err);
      setDeleteError('Falha ao excluir o evento. Tente novamente.');
    } finally {
      setIsDeletingEvent(false);
    }
  };

  const handleCreateEvent = async () => {
    try {
      setIsCreatingEvent(true);
      console.log('Creating event with data:', {
        eventData: newEventData,
        images: newEventImages,
      });

      // Convert images to base64
      const imagesBase64 = await api.filesToBase64(newEventImages);

      // Convert date to ISO format
      const isoDate = new Date(newEventData.date).toISOString();

      // Create the galery event
      const createdEvent = await api.createGaleryEvent({
        name: newEventData.name,
        location: newEventData.location,
        date: isoDate,
        images_base64: imagesBase64,
      });

      console.log('Event created successfully:', createdEvent);

      // Refresh the gallery data
      await fetchGalleryData();

      // Show success toast
      toast({
        title: "Sucesso!",
        description: "Evento criado com sucesso!",
      });

      // Reset form
      setNewEventData({ name: "", date: "", location: "", description: "" });
      setNewEventImages([]);
      setShowCreateEvent(false);
    } catch (err) {
      console.error('Failed to create event:', err);
      toast({
        variant: "destructive",
        title: "Erro ao criar evento",
        description: "Falha ao criar evento. Verifique o console para mais detalhes.",
      });
    } finally {
      setIsCreatingEvent(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-background py-12">
        <div className="max-w-7xl mx-auto px-4">
          <div className="text-center">
            <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-6">
              Galeria de Fotos
            </h1>
            <p className="text-xl text-muted-foreground">Carregando...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background py-12">
        <div className="max-w-7xl mx-auto px-4">
          <div className="text-center">
            <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-6">
              Galeria de Fotos
            </h1>
            <div className="text-destructive mb-4">{error}</div>
            <Button onClick={fetchGalleryData}>Tentar Novamente</Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background py-12">
      <div className="max-w-7xl mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-6">
            Galeria de Fotos
          </h1>
          <p className="text-xl text-muted-foreground max-w-3xl mx-auto">
            Reviva os melhores momentos dos nossos eventos e veja como nossa comunidade cresce a cada encontro.
          </p>
        </div>

        {/* Events Section */}
        <div className="mb-16">
          <div className="flex items-center justify-between mb-8">
            <h2 className="text-3xl font-bold text-foreground">Eventos</h2>
            <ShowWhenAuthenticated>
              <Button onClick={() => setShowCreateEvent(true)} className="gap-2">
                <Plus className="h-4 w-4" />
                Criar Evento
              </Button>
            </ShowWhenAuthenticated>
          </div>

          {events.length === 0 ? (
            <div className="text-center py-16 bg-muted rounded-lg">
              <ImageIcon className="h-16 w-16 mx-auto text-muted-foreground/50 mb-4" />
              <h3 className="text-lg font-medium text-foreground mb-2">Nenhum evento ainda</h3>
              <p className="text-muted-foreground mb-4">
                Crie o primeiro evento e compartilhe fotos da comunidade!
              </p>
              <ShowWhenAuthenticated>
                <Button onClick={() => setShowCreateEvent(true)}>
                  Criar Primeiro Evento
                </Button>
              </ShowWhenAuthenticated>
            </div>
          ) : (
            <div className="grid lg:grid-cols-2 gap-8">
              {events.map((event) => (
                <Card key={event.id} className="hover:shadow-lg transition-shadow">
                  <CardHeader>
                    <CardTitle className="flex items-start justify-between gap-4">
                      <span>{event.name}</span>
                      <div className="flex items-center gap-2 shrink-0">
                        <Badge variant="secondary">
                          {event.photos.length} fotos
                        </Badge>
                        <ShowWhenAuthenticated>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
                            onClick={() => setEventToDelete(event)}
                          >
                            <X className="h-4 w-4" />
                          </Button>
                        </ShowWhenAuthenticated>
                      </div>
                    </CardTitle>
                    <div className="flex flex-wrap gap-3 text-sm text-muted-foreground">
                      <div className="flex items-center gap-1">
                        <Calendar className="h-3 w-3" />
                        {formatDate(event.date)}
                      </div>
                      <div className="flex items-center gap-1">
                        <MapPin className="h-3 w-3" />
                        {event.location}
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent>
                    {event.photos.length > 0 ? (
                      <div className="grid grid-cols-3 gap-2">
                        {event.photos.slice(0, 6).map((photo, index) => (
                          <div
                            key={`${event.id}-photo-${index}`}
                            className="aspect-square bg-muted rounded-lg cursor-pointer hover:opacity-80 transition-opacity overflow-hidden"
                            onClick={() =>
                              setSelectedPhoto({
                                ...photo,
                                eventName: event.name,
                                eventDate: event.date,
                                eventLocation: event.location,
                                source: "event",
                                eventId: event.id,
                              })
                            }
                          >
                            <img
                              src={photo.url}
                              alt={photo.caption}
                              className="w-full h-full object-cover"
                              onError={(e) => {
                                // Fallback if image fails to load
                                (e.target as HTMLImageElement).style.display = 'none';
                                const parent = (e.target as HTMLImageElement).parentElement;
                                if (parent) {
                                  parent.innerHTML = '<div class="w-full h-full bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center"><svg class="h-8 w-8 text-muted-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg></div>';
                                }
                              }}
                            />
                          </div>
                        ))}
                        {event.photos.length > 6 && (
                          <div className="aspect-square bg-muted rounded-lg flex items-center justify-center">
                            <span className="text-sm text-muted-foreground">
                              +{event.photos.length - 6} fotos
                            </span>
                          </div>
                        )}
                      </div>
                    ) : (
                      <div className="text-center py-8 bg-muted rounded-lg">
                        <ImageIcon className="h-12 w-12 mx-auto text-muted-foreground/50 mb-2" />
                        <p className="text-muted-foreground">Nenhuma foto ainda</p>
                      </div>
                    )}
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>

        {/* All Photos Gallery */}
        <div>
          <div className="flex items-center justify-between mb-8">
            <h2 className="text-3xl font-bold text-foreground">Mais fotos</h2>
            <ShowWhenAuthenticated>
              <Button onClick={() => setShowImageUpload(true)} className="gap-2">
                <ImagePlus className="h-4 w-4" />
                Adicione Imagens
              </Button>
            </ShowWhenAuthenticated>
          </div>

          {allPhotos.length === 0 ? (
            <p className="text-center text-muted-foreground py-8">Nenhuma imagem adicionada ainda</p>
          ) : (
            <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              {allPhotos.map((photo) => (
                <div
                  key={photo.id}
                  className="relative aspect-square bg-muted rounded-lg cursor-pointer hover:opacity-80 transition-opacity overflow-hidden group"
                  onClick={() =>
                    setSelectedPhoto({
                      ...photo,
                      source: "gallery",
                    })
                  }
                >
                  <img
                    src={photo.url}
                    alt={photo.caption}
                    className="w-full h-full object-cover"
                    onError={(e) => {
                      // Fallback if image fails to load
                      const target = e.target as HTMLImageElement;
                      target.style.display = 'none';
                      const parent = target.parentElement;
                      if (parent) {
                        parent.innerHTML = '<div class="w-full h-full bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center"><svg class="h-12 w-12 text-muted-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg></div>';
                      }
                    }}
                  />
                  <div className="absolute inset-x-0 bottom-0 bg-black/60 text-white p-2 transform translate-y-full group-hover:translate-y-0 transition-transform">
                    <p className="text-xs font-medium truncate">{photo.eventName || photo.caption}</p>
                    {photo.eventDate && <p className="text-xs opacity-75">{formatDate(photo.eventDate)}</p>}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Photo Modal */}
        {selectedPhoto && (
          <Dialog open={!!selectedPhoto} onOpenChange={() => setSelectedPhoto(null)}>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle>{selectedPhoto.caption}</DialogTitle>
              </DialogHeader>

              {deleteError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{deleteError}</AlertDescription>
                </Alert>
              )}

              <div className="space-y-4">
                <div className="aspect-video bg-muted rounded-lg overflow-hidden">
                  <img
                    src={selectedPhoto.url}
                    alt={selectedPhoto.caption}
                    className="w-full h-full object-contain"
                    onError={(e) => {
                      const target = e.target as HTMLImageElement;
                      target.style.display = 'none';
                      const parent = target.parentElement;
                      if (parent) {
                        parent.innerHTML = '<div class="w-full h-full flex items-center justify-center"><svg class="h-16 w-16 text-muted-foreground" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg></div>';
                      }
                    }}
                  />
                </div>

                {(selectedPhoto.eventDate ||
                  selectedPhoto.eventLocation ||
                  selectedPhoto.eventName) && (
                    <div className="space-y-2">
                      {selectedPhoto.eventDate && (
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <Calendar className="h-3 w-3" />
                          {formatDate(selectedPhoto.eventDate)}
                        </div>
                      )}
                      {selectedPhoto.eventLocation && (
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <MapPin className="h-3 w-3" />
                          {selectedPhoto.eventLocation}
                        </div>
                      )}
                      {selectedPhoto.eventName && (
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <Users className="h-3 w-3" />
                          Evento: {selectedPhoto.eventName}
                        </div>
                      )}
                    </div>
                  )}
              </div>

              {/* Only show delete button for images with real Firestore IDs */}
              <ShowWhenAuthenticated>
                <DialogFooter>
                  {/* Gallery photo - delete the image resource */}
                  {selectedPhoto.source === "gallery" && canDeletePhoto(selectedPhoto.id) && (
                    <Button
                      variant="destructive"
                      onClick={handleDeleteImage}
                      disabled={isDeleting}
                      className="gap-2"
                    >
                      <Trash2 className="h-4 w-4" />
                      {isDeleting ? "Excluindo..." : "Excluir Imagem"}
                    </Button>
                  )}

                  {/* Event photo - remove the photo from the event via modifyGaleryEventRequest */}
                  {selectedPhoto.source === "event" && (
                    <Button
                      variant="destructive"
                      onClick={() => handleDeleteEventPhoto(selectedPhoto)}
                      disabled={isDeleting}
                      className="gap-2"
                    >
                      <Trash2 className="h-4 w-4" />
                      {isDeleting ? "Excluindo..." : "Remover do evento"}
                    </Button>
                  )}
                </DialogFooter>
              </ShowWhenAuthenticated>
            </DialogContent>
          </Dialog>
        )}

        {/* Delete Event Confirmation Modal */}
        {eventToDelete && (
          <Dialog open={!!eventToDelete} onOpenChange={() => setEventToDelete(null)}>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle>Excluir Evento</DialogTitle>
              </DialogHeader>

              {deleteError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{deleteError}</AlertDescription>
                </Alert>
              )}

              <div className="space-y-4">
                <p className="text-muted-foreground">
                  Tem certeza que deseja excluir o evento <strong>{eventToDelete.name}</strong>?
                </p>
                <p className="text-sm text-muted-foreground">
                  Esta ação não pode ser desfeita. O evento será permanentemente removido.
                </p>
              </div>

              <DialogFooter className="gap-2 sm:gap-0">
                <Button
                  variant="outline"
                  onClick={() => setEventToDelete(null)}
                  disabled={isDeletingEvent}
                >
                  Cancelar
                </Button>
                <Button
                  variant="destructive"
                  onClick={handleDeleteEvent}
                  disabled={isDeletingEvent}
                  className="gap-2"
                >
                  <Trash2 className="h-4 w-4" />
                  {isDeletingEvent ? 'Excluindo...' : 'Excluir Evento'}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        )}

        {/* Image Upload Modal */}
        <FileUploadModal
          open={showImageUpload}
          onOpenChange={setShowImageUpload}
          onUpload={handleUploadImage}
          title="Adicionar Imagens"
          uploadButtonText="Enviar"
          config={{
            accept: "image/*",
            maxSize: 10 * 1024 * 1024, // 10MB
            multiple: true,
            fileCategory: "image",
          }}
        />

        {/* Create Event Modal */}
        <CreateEventModal
          open={showCreateEvent}
          onOpenChange={setShowCreateEvent}
          eventData={newEventData}
          onEventDataChange={setNewEventData}
          eventImages={newEventImages}
          onEventImagesChange={setNewEventImages}
          onSubmit={handleCreateEvent}
          isCreatingEvent={isCreatingEvent}
        />
      </div>
    </div>
  );
};

// Create Event Modal Component
interface CreateEventModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  eventData: {
    name: string;
    date: string;
    location: string;
    description: string;
  };
  onEventDataChange: (data: any) => void;
  eventImages: File[];
  onEventImagesChange: (images: File[]) => void;
  onSubmit: () => void;
  isCreatingEvent: boolean;
}

const CreateEventModal: React.FC<CreateEventModalProps> = ({
  open,
  onOpenChange,
  eventData,
  onEventDataChange,
  eventImages,
  onEventImagesChange,
  onSubmit,
  isCreatingEvent,
}) => {
  const [showImageUpload, setShowImageUpload] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit();
  };

  const handleClose = () => {
    onOpenChange(false);
  };

  const removeImage = (index: number) => {
    onEventImagesChange(eventImages.filter((_, i) => i !== index));
  };

  const isFormValid =
    eventData.name.trim() !== "" &&
    eventData.date !== "" &&
    eventData.location.trim() !== "" &&
    eventImages.length > 0; // Require at least one image

  return (
    <>
      <Dialog open={open} onOpenChange={handleClose}>
        <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Plus className="h-5 w-5" />
              Criar Novo Evento
            </DialogTitle>
          </DialogHeader>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Event Details */}
            <div className="space-y-4">
              <div>
                <Label htmlFor="event-name">Nome do Evento *</Label>
                <Input
                  id="event-name"
                  value={eventData.name}
                  onChange={(e) =>
                    onEventDataChange({ ...eventData, name: e.target.value })
                  }
                  placeholder="Ex: Workshop Python Básico"
                  required
                />
              </div>

              <div>
                <Label htmlFor="event-date">Data *</Label>
                <Input
                  id="event-date"
                  type="date"
                  value={eventData.date}
                  onChange={(e) =>
                    onEventDataChange({ ...eventData, date: e.target.value })
                  }
                  required
                />
              </div>

              <div>
                <Label htmlFor="event-location">Localização *</Label>
                <Input
                  id="event-location"
                  value={eventData.location}
                  onChange={(e) =>
                    onEventDataChange({ ...eventData, location: e.target.value })
                  }
                  placeholder="Ex: IFSP São Carlos"
                  required
                />
              </div>

              <div>
                <Label htmlFor="event-description">Descrição do Evento</Label>
                <Textarea
                  id="event-description"
                  value={eventData.description}
                  onChange={(e) =>
                    onEventDataChange({ ...eventData, description: e.target.value })
                  }
                  placeholder="Breve descrição do evento..."
                  rows={4}
                />
              </div>
            </div>

            {/* Event Images */}
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Imagens do Evento *</Label>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => setShowImageUpload(true)}
                  className="gap-2"
                >
                  <ImagePlus className="h-4 w-4" />
                  Adicionar Imagens
                </Button>
              </div>

              {eventImages.length > 0 && (
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">
                    {eventImages.length} imagem(ns) selecionada(s)
                  </p>
                  <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                    {eventImages.map((image, index) => (
                      <div
                        key={index}
                        className="relative group aspect-square bg-muted rounded-lg overflow-hidden"
                      >
                        <img
                          src={URL.createObjectURL(image)}
                          alt={image.name}
                          className="w-full h-full object-cover"
                        />
                        <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                          <Button
                            type="button"
                            variant="destructive"
                            size="sm"
                            onClick={() => removeImage(index)}
                          >
                            <X className="h-4 w-4" />
                          </Button>
                        </div>
                        <div className="absolute bottom-0 left-0 right-0 bg-black/60 text-white p-1">
                          <p className="text-xs truncate">{image.name}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {eventImages.length === 0 && (
                <div className="text-center py-8 border-2 border-dashed rounded-lg border-muted-foreground/25">
                  <ImageIcon className="h-12 w-12 mx-auto text-muted-foreground/50 mb-2" />
                  <p className="text-sm text-muted-foreground">
                    Nenhuma imagem adicionada ainda
                  </p>
                  <p className="text-xs text-muted-foreground mt-1">
                    Pelo menos uma imagem é obrigatória
                  </p>
                </div>
              )}
            </div>

            {/* Action Buttons */}
            <div className="flex gap-2 justify-end pt-4 border-t">
              <Button type="button" variant="outline" onClick={handleClose}>
                Cancelar
              </Button>
              <Button type="submit" disabled={!isFormValid || isCreatingEvent} className="gap-2">
                {isCreatingEvent ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Plus className="h-4 w-4" />
                )}
                {isCreatingEvent ? "Criando..." : "Criar Evento"}
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Nested Image Upload Modal */}
      <FileUploadModal
        open={showImageUpload}
        onOpenChange={setShowImageUpload}
        onUpload={(files) => {
          onEventImagesChange([...eventImages, ...files]);
          setShowImageUpload(false);
        }}
        title="Adicionar Imagens ao Evento"
        uploadButtonText="Adicionar"
        config={{
          accept: "image/*",
          maxSize: 10 * 1024 * 1024, // 10MB
          multiple: true,
          fileCategory: "image",
        }}
      />
    </>
  );
};

export default Galeria;
