import { useState } from "react";
import { Calendar, MapPin, Users, Award, Plus, Image as ImageIcon, X } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { ShowWhenAuthenticated } from "@/auth/AuthSwitch";
import { FileUploadModal } from "@/components/FileUploadModal";

interface TimelineEvent {
  id: string;
  date: string;
  title: string;
  description: string;
  location: string;
  images?: string[];
}

const Historia = () => {
  const [timelineEvents, setTimelineEvents] = useState<TimelineEvent[]>([
    {
      id: "1",
      date: "2015-03-15",
      title: "Fundação do Grupy Sanca",
      description: "O grupo nasceu da necessidade de criar uma comunidade local de Python em São Carlos, inspirado por outros grupys pelo Brasil.",
      location: "IFSP São Carlos"
    },
    {
      id: "2",
      date: "2016-05-20",
      title: "Primeiro Workshop Oficial",
      description: "Organizamos nosso primeiro workshop sobre Django, marcando o início das atividades educacionais regulares.",
      location: "USP São Carlos"
    },
    {
      id: "3",
      date: "2017-08-10",
      title: "Parcerias com Universidades",
      description: "Estabelecemos parcerias formais com USP e IFSP para sediar eventos e atingir mais estudantes.",
      location: "Múltiplas instituições"
    },
    {
      id: "4",
      date: "2018-11-18",
      title: "Primeiro Python Day São Carlos",
      description: "Organizamos um evento de dia inteiro com palestras, workshops e networking, nosso maior evento até então.",
      location: "Centro de Convenções"
    },
    {
      id: "5",
      date: "2019-06-22",
      title: "Expansão Regional",
      description: "O grupo começou a receber participantes de cidades vizinhas, consolidando-se como referência regional.",
      location: "São Carlos e região"
    },
    {
      id: "6",
      date: "2020-04-01",
      title: "Eventos Online",
      description: "Durante a pandemia, adaptamos todos os eventos para formato online, mantendo a comunidade ativa e unida.",
      location: "Online"
    },
    {
      id: "7",
      date: "2022-09-15",
      title: "Retorno Híbrido",
      description: "Retomamos os eventos presenciais com transmissão online, ampliando nosso alcance e inclusividade.",
      location: "Híbrido"
    },
    {
      id: "8",
      date: "2023-03-15",
      title: "8 Anos de Comunidade",
      description: "Celebramos 8 anos de atividades contínuas, com mais de 50 eventos realizados e centenas de pessoas impactadas.",
      location: "São Carlos"
    }
  ]);

  const [showCreateTimeline, setShowCreateTimeline] = useState(false);
  const [newTimelineData, setNewTimelineData] = useState({
    date: "",
    title: "",
    description: "",
    location: "",
  });
  const [newTimelineImages, setNewTimelineImages] = useState<File[]>([]);

  const achievements = [
    {
      icon: Calendar,
      title: "50+ Eventos",
      description: "Workshops, palestras e meetups realizados"
    },
    {
      icon: Users,
      title: "500+ Participantes",
      description: "Pessoas diferentes já participaram de nossos eventos"
    },
    {
      icon: Award,
      title: "Reconhecimento Nacional",
      description: "Participação ativa na comunidade Python Brasil"
    },
    {
      icon: MapPin,
      title: "Impacto Regional",
      description: "Referência para São Carlos e região"
    }
  ];

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('pt-BR');
  };

  const getYearFromDate = (dateString: string) => {
    return new Date(dateString).getFullYear().toString();
  };

  const handleCreateTimeline = () => {
    console.log('CreateTimelineEntry triggered with data:', {
      timelineData: newTimelineData,
      images: newTimelineImages,
    });
    // TODO: Implement actual timeline entry creation logic

    // Reset form
    setNewTimelineData({ date: "", title: "", description: "", location: "" });
    setNewTimelineImages([]);
    setShowCreateTimeline(false);
  };

  return (
    <div className="min-h-screen bg-background py-12">
      <div className="max-w-6xl mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-16">
          <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-6">
            Nossa História
          </h1>
          <p className="text-xl text-muted-foreground max-w-3xl mx-auto leading-relaxed">
            Conheça a trajetória do Grupy Sanca, desde sua fundação até os dias atuais,
            e como nos tornamos uma referência na comunidade Python regional.
          </p>
        </div>

        {/* Origin Story */}
        <div className="gradient-section p-8 rounded-xl mb-16">
          <h2 className="text-3xl font-bold text-foreground mb-6">Como tudo começou</h2>
          <div className="prose prose-lg text-muted-foreground leading-relaxed">
            <p className="mb-4">
              Em 2015, um pequeno grupo de entusiastas do Python em São Carlos percebeu a necessidade
              de criar um espaço local para compartilhar conhecimento e experiências sobre esta linguagem
              que estava crescendo rapidamente no cenário tecnológico brasileiro.
            </p>
            <p className="mb-4">
              Inspirados pelos grupos de usuários Python (Grupys) que já existiam em outras cidades como
              São Paulo, Rio de Janeiro e Campinas, decidimos fundar o Grupy Sanca com o objetivo
              de democratizar o acesso ao conhecimento sobre Python na região.
            </p>
            <p>
              O que começou como encontros informais em cafés e salas emprestadas, rapidamente evoluiu
              para uma comunidade estruturada que hoje é referência no interior paulista.
            </p>
          </div>
        </div>

        {/* Achievements */}
        <div className="mb-16">
          <h2 className="text-3xl font-bold text-foreground mb-8 text-center">
            Nossos números
          </h2>
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
            {achievements.map((achievement, index) => (
              <Card key={index} className="text-center hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="mx-auto bg-primary/10 p-3 rounded-full w-fit mb-4">
                    <achievement.icon className="h-8 w-8 text-primary" />
                  </div>
                  <CardTitle className="text-lg">{achievement.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-muted-foreground text-sm">
                    {achievement.description}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>

        {/* Timeline */}
        <div className="mb-16">
          <div className="flex items-center justify-between mb-8">
            <h2 className="text-3xl font-bold text-foreground">
              Linha do tempo
            </h2>
            <ShowWhenAuthenticated>
              <Button onClick={() => setShowCreateTimeline(true)} className="gap-2">
                <Plus className="h-4 w-4" />
                Adicionar Entrada
              </Button>
            </ShowWhenAuthenticated>
          </div>
          <div className="relative">
            {/* Continuous timeline line */}
            <div
              className="absolute left-[30px] top-[30px] w-1 bg-primary"
              style={{ height: `calc(100% - 60px)` }}
            ></div>

            {timelineEvents.map((event, index) => (
              <div key={event.id} className="flex gap-6 items-start mb-8 last:mb-0">
                <div className="flex flex-col items-center relative w-[60px]">
                  <div className="bg-primary text-primary-foreground rounded-full p-3 font-bold text-lg w-[60px] h-[60px] flex items-center justify-center relative z-10">
                    {getYearFromDate(event.date)}
                  </div>
                </div>

                <Card className="flex-1 hover:shadow-lg transition-shadow">
                  <CardHeader>
                    <CardTitle className="flex items-start justify-between flex-wrap gap-2">
                      <div className="flex flex-col gap-1">
                        <span>{event.title}</span>
                        <div className="flex items-center gap-1 text-sm font-normal text-muted-foreground">
                          <Calendar className="h-3 w-3" />
                          {formatDate(event.date)}
                        </div>
                      </div>
                      <Badge variant="outline" className="flex items-center gap-1">
                        <MapPin className="h-3 w-3" />
                        {event.location}
                      </Badge>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      {event.description}
                    </p>
                    {/* Image placeholders - exclude first entry and specific years */}
                    {index !== 0 && getYearFromDate(event.date) !== "2020" && (
                      <div className={`grid gap-3 mt-4 ${index === 1 || index === 3 || index === 6 ? 'grid-cols-2' : 'grid-cols-1'}`}>
                        {((index === 1 || index === 3 || index === 6) ? [1, 2] : [1]).map((imgIndex) => (
                          <div key={imgIndex} className="bg-muted rounded-lg h-32 flex items-center justify-center">
                            <div className="text-center text-muted-foreground">
                              <Calendar className="h-8 w-8 mx-auto mb-1 opacity-50" />
                              <p className="text-xs">Foto {imgIndex}</p>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </CardContent>
                </Card>
              </div>
            ))}
          </div>
        </div>

        {/* Create Timeline Entry Modal */}
        <CreateTimelineEntryModal
          open={showCreateTimeline}
          onOpenChange={setShowCreateTimeline}
          timelineData={newTimelineData}
          onTimelineDataChange={setNewTimelineData}
          timelineImages={newTimelineImages}
          onTimelineImagesChange={setNewTimelineImages}
          onSubmit={handleCreateTimeline}
        />
      </div>
    </div>
  );
};

// Create Timeline Entry Modal Component
interface CreateTimelineEntryModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  timelineData: {
    date: string;
    title: string;
    description: string;
    location: string;
  };
  onTimelineDataChange: (data: any) => void;
  timelineImages: File[];
  onTimelineImagesChange: (images: File[]) => void;
  onSubmit: () => void;
}

const CreateTimelineEntryModal: React.FC<CreateTimelineEntryModalProps> = ({
  open,
  onOpenChange,
  timelineData,
  onTimelineDataChange,
  timelineImages,
  onTimelineImagesChange,
  onSubmit,
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
    onTimelineImagesChange(timelineImages.filter((_, i) => i !== index));
  };

  const isFormValid =
    timelineData.date !== "" &&
    timelineData.title.trim() !== "" &&
    timelineData.description.trim() !== "" &&
    timelineData.location.trim() !== "";

  return (
    <>
      <Dialog open={open} onOpenChange={handleClose}>
        <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Plus className="h-5 w-5" />
              Adicionar Entrada na Linha do Tempo
            </DialogTitle>
          </DialogHeader>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Timeline Entry Details */}
            <div className="space-y-4">
              <div>
                <Label htmlFor="timeline-title">Nome/Título *</Label>
                <Input
                  id="timeline-title"
                  value={timelineData.title}
                  onChange={(e) =>
                    onTimelineDataChange({ ...timelineData, title: e.target.value })
                  }
                  placeholder="Ex: Primeiro Hackathon Regional"
                  required
                />
              </div>

              <div>
                <Label htmlFor="timeline-date">Data *</Label>
                <Input
                  id="timeline-date"
                  type="date"
                  value={timelineData.date}
                  onChange={(e) =>
                    onTimelineDataChange({ ...timelineData, date: e.target.value })
                  }
                  required
                />
              </div>

              <div>
                <Label htmlFor="timeline-location">Localização *</Label>
                <Input
                  id="timeline-location"
                  value={timelineData.location}
                  onChange={(e) =>
                    onTimelineDataChange({ ...timelineData, location: e.target.value })
                  }
                  placeholder="Ex: IFSP São Carlos"
                  required
                />
              </div>

              <div>
                <Label htmlFor="timeline-description">Descrição *</Label>
                <Textarea
                  id="timeline-description"
                  value={timelineData.description}
                  onChange={(e) =>
                    onTimelineDataChange({ ...timelineData, description: e.target.value })
                  }
                  placeholder="Descreva este momento importante na história..."
                  rows={4}
                  required
                />
              </div>
            </div>

            {/* Timeline Images */}
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Imagens</Label>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => setShowImageUpload(true)}
                  className="gap-2"
                >
                  <ImageIcon className="h-4 w-4" />
                  Adicionar Imagens
                </Button>
              </div>

              {timelineImages.length > 0 && (
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">
                    {timelineImages.length} imagem(ns) selecionada(s)
                  </p>
                  <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                    {timelineImages.map((image, index) => (
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

              {timelineImages.length === 0 && (
                <div className="text-center py-8 border-2 border-dashed rounded-lg border-muted-foreground/25">
                  <ImageIcon className="h-12 w-12 mx-auto text-muted-foreground/50 mb-2" />
                  <p className="text-sm text-muted-foreground">
                    Nenhuma imagem adicionada ainda
                  </p>
                </div>
              )}
            </div>

            {/* Action Buttons */}
            <div className="flex gap-2 justify-end pt-4 border-t">
              <Button type="button" variant="outline" onClick={handleClose}>
                Cancelar
              </Button>
              <Button type="submit" disabled={!isFormValid} className="gap-2">
                <Plus className="h-4 w-4" />
                Adicionar Entrada
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
          onTimelineImagesChange([...timelineImages, ...files]);
          setShowImageUpload(false);
        }}
        title="Adicionar Imagens à Entrada"
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

export default Historia;
