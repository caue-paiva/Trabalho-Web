import { useState } from "react";
import { Upload, Plus, Calendar, MapPin, Users, Image as ImageIcon, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";

interface Event {
  id: string;
  name: string;
  date: string;
  location: string;
  description: string;
  photos: Photo[];
}

interface Photo {
  id: string;
  url: string;
  caption: string;
  eventId: string;
  uploadedAt: string;
}

const Galeria = () => {
  const [events, setEvents] = useState<Event[]>([
    {
      id: "1",
      name: "Python Day São Carlos 2023",
      date: "2023-11-18",
      location: "IFSP São Carlos",
      description: "Nosso maior evento anual com palestras, workshops e networking.",
      photos: [
        {
          id: "1",
          url: "/api/placeholder/400/300",
          caption: "Abertura do evento",
          eventId: "1",
          uploadedAt: "2023-11-18T09:00:00"
        },
        {
          id: "2", 
          url: "/api/placeholder/400/300",
          caption: "Workshop de Django",
          eventId: "1",
          uploadedAt: "2023-11-18T14:00:00"
        },
        {
          id: "3",
          url: "/api/placeholder/400/300", 
          caption: "Networking no coffee break",
          eventId: "1",
          uploadedAt: "2023-11-18T15:30:00"
        }
      ]
    },
    {
      id: "2",
      name: "Workshop Flask & APIs",
      date: "2023-10-14",
      location: "USP São Carlos",
      description: "Workshop prático sobre desenvolvimento de APIs com Flask.",
      photos: [
        {
          id: "4",
          url: "/api/placeholder/400/300",
          caption: "Explicação sobre REST APIs",
          eventId: "2", 
          uploadedAt: "2023-10-14T10:00:00"
        },
        {
          id: "5",
          url: "/api/placeholder/400/300",
          caption: "Hands-on coding session",
          eventId: "2",
          uploadedAt: "2023-10-14T11:30:00"
        }
      ]
    }
  ]);

  const [selectedPhoto, setSelectedPhoto] = useState<Photo | null>(null);
  const [showUploadForm, setShowUploadForm] = useState(false);
  const [newEvent, setNewEvent] = useState({
    name: "",
    date: "",
    location: "",
    description: ""
  });

  const allPhotos = events.flatMap(event => 
    event.photos.map(photo => ({
      ...photo,
      eventName: event.name,
      eventDate: event.date,
      eventLocation: event.location
    }))
  ).sort((a, b) => new Date(b.uploadedAt).getTime() - new Date(a.uploadedAt).getTime());

  const createEvent = () => {
    const event: Event = {
      id: Date.now().toString(),
      ...newEvent,
      photos: []
    };
    setEvents([...events, event]);
    setNewEvent({ name: "", date: "", location: "", description: "" });
    setShowUploadForm(false);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('pt-BR');
  };

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

        {/* Upload Section */}
        <div className="mb-12">
          <Card className="gradient-section">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Upload className="h-5 w-5" />
                Contribua com fotos
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground mb-4">
                Participou de algum evento? Ajude a documentar nossa história compartilhando suas fotos!
              </p>
              
              <Dialog open={showUploadForm} onOpenChange={setShowUploadForm}>
                <DialogTrigger asChild>
                  <Button className="gap-2">
                    <Plus className="h-4 w-4" />
                    Criar evento e adicionar fotos
                  </Button>
                </DialogTrigger>
                <DialogContent className="sm:max-w-md">
                  <DialogHeader>
                    <DialogTitle>Criar novo evento</DialogTitle>
                    <DialogDescription>
                      Primeiro, vamos criar um evento para organizar suas fotos.
                    </DialogDescription>
                  </DialogHeader>
                  <div className="space-y-4">
                    <div>
                      <Label htmlFor="name">Nome do evento</Label>
                      <Input
                        id="name"
                        value={newEvent.name}
                        onChange={(e) => setNewEvent({ ...newEvent, name: e.target.value })}
                        placeholder="Ex: Workshop Python Básico"
                      />
                    </div>
                    <div>
                      <Label htmlFor="date">Data</Label>
                      <Input
                        id="date"
                        type="date"
                        value={newEvent.date}
                        onChange={(e) => setNewEvent({ ...newEvent, date: e.target.value })}
                      />
                    </div>
                    <div>
                      <Label htmlFor="location">Local</Label>
                      <Input
                        id="location"
                        value={newEvent.location}
                        onChange={(e) => setNewEvent({ ...newEvent, location: e.target.value })}
                        placeholder="Ex: IFSP São Carlos"
                      />
                    </div>
                    <div>
                      <Label htmlFor="description">Descrição</Label>
                      <Textarea
                        id="description"
                        value={newEvent.description}
                        onChange={(e) => setNewEvent({ ...newEvent, description: e.target.value })}
                        placeholder="Breve descrição do evento..."
                        rows={3}
                      />
                    </div>
                    <Button 
                      onClick={createEvent} 
                      className="w-full"
                      disabled={!newEvent.name || !newEvent.date || !newEvent.location}
                    >
                      Criar evento
                    </Button>
                  </div>
                </DialogContent>
              </Dialog>
            </CardContent>
          </Card>
        </div>

        {/* Events Section */}
        <div className="mb-16">
          <h2 className="text-3xl font-bold text-foreground mb-8">Eventos recentes</h2>
          <div className="grid lg:grid-cols-2 gap-8">
            {events.map((event) => (
              <Card key={event.id} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <CardTitle className="flex items-start justify-between gap-4">
                    <span>{event.name}</span>
                    <Badge variant="secondary" className="shrink-0">
                      {event.photos.length} fotos
                    </Badge>
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
                  <p className="text-muted-foreground mb-4">{event.description}</p>
                  
                  {event.photos.length > 0 ? (
                    <div className="grid grid-cols-3 gap-2">
                      {event.photos.slice(0, 6).map((photo) => (
                        <div
                          key={photo.id}
                          className="aspect-square bg-muted rounded-lg cursor-pointer hover:opacity-80 transition-opacity overflow-hidden"
                          onClick={() => setSelectedPhoto({
                            ...photo,
                            eventName: event.name,
                            eventDate: event.date,
                            eventLocation: event.location
                          } as any)}
                        >
                          <div className="w-full h-full bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center">
                            <ImageIcon className="h-8 w-8 text-muted-foreground" />
                          </div>
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
                      <Button variant="outline" size="sm" className="mt-2">
                        Adicionar fotos
                      </Button>
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        </div>

        {/* All Photos Gallery */}
        <div>
          <h2 className="text-3xl font-bold text-foreground mb-8">Todas as fotos</h2>
          <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {allPhotos.map((photo) => (
              <div
                key={photo.id}
                className="aspect-square bg-muted rounded-lg cursor-pointer hover:opacity-80 transition-opacity overflow-hidden group"
                onClick={() => setSelectedPhoto(photo)}
              >
                <div className="w-full h-full bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center relative">
                  <ImageIcon className="h-12 w-12 text-muted-foreground group-hover:scale-110 transition-transform" />
                  <div className="absolute inset-x-0 bottom-0 bg-black/60 text-white p-2 transform translate-y-full group-hover:translate-y-0 transition-transform">
                    <p className="text-xs font-medium truncate">{photo.eventName}</p>
                    <p className="text-xs opacity-75">{formatDate(photo.eventDate)}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
          
          {allPhotos.length === 0 && (
            <div className="text-center py-16 bg-muted rounded-lg">
              <ImageIcon className="h-16 w-16 mx-auto text-muted-foreground/50 mb-4" />
              <h3 className="text-lg font-medium text-foreground mb-2">Ainda não temos fotos</h3>
              <p className="text-muted-foreground mb-4">
                Seja o primeiro a compartilhar fotos de nossos eventos!
              </p>
              <Button onClick={() => setShowUploadForm(true)}>
                Adicionar primeira foto
              </Button>
            </div>
          )}
        </div>

        {/* Photo Modal */}
        {selectedPhoto && (
          <Dialog open={!!selectedPhoto} onOpenChange={() => setSelectedPhoto(null)}>
            <DialogContent className="sm:max-w-lg">
              <DialogHeader>
                <DialogTitle className="flex items-center justify-between">
                  <span>{selectedPhoto.caption}</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSelectedPhoto(null)}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </DialogTitle>
              </DialogHeader>
              
              <div className="space-y-4">
                <div className="aspect-video bg-muted rounded-lg flex items-center justify-center">
                  <ImageIcon className="h-16 w-16 text-muted-foreground" />
                </div>
                
                <div className="space-y-2">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Calendar className="h-3 w-3" />
                    {formatDate((selectedPhoto as any).eventDate)}
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <MapPin className="h-3 w-3" />
                    {(selectedPhoto as any).eventLocation}
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Users className="h-3 w-3" />
                    Evento: {(selectedPhoto as any).eventName}
                  </div>
                </div>
              </div>
            </DialogContent>
          </Dialog>
        )}
      </div>
    </div>
  );
};

export default Galeria;