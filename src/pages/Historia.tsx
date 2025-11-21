import { useState, useEffect } from "react";
import { Calendar, MapPin, Users, Award, Plus, X, Trash2, AlertCircle } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ShowWhenAuthenticated } from "@/auth/AuthSwitch";
import * as api from "@/services/api";

interface TimelineEvent {
  id: string;
  date: string;
  title: string;
  description: string;
  location: string;
}

const Historia = () => {
  const [timelineEvents, setTimelineEvents] = useState<TimelineEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [showCreateTimeline, setShowCreateTimeline] = useState(false);
  const [newTimelineData, setNewTimelineData] = useState({
    date: "",
    title: "",
    description: "",
    location: "",
  });

  const [entryToDelete, setEntryToDelete] = useState<TimelineEvent | null>(null);
  const [isDeletingEntry, setIsDeletingEntry] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  // Fetch timeline entries on mount
  useEffect(() => {
    fetchTimelineEntries();
  }, []);

  const fetchTimelineEntries = async () => {
    try {
      setLoading(true);
      setError(null);

      const entries = await api.listTimelineEntries();

      // Transform API response to UI format
      const transformedEntries: TimelineEvent[] = entries.map(entry => ({
        id: entry.id,
        date: entry.date,
        title: entry.name,
        description: entry.text,
        location: entry.location,
      }));

      // Sort by date descending (newest first)
      transformedEntries.sort((a, b) =>
        new Date(b.date).getTime() - new Date(a.date).getTime()
      );

      setTimelineEvents(transformedEntries);
    } catch (err) {
      console.error('Failed to fetch timeline entries:', err);
      setError('Falha ao carregar entradas da linha do tempo. Tente novamente mais tarde.');
    } finally {
      setLoading(false);
    }
  };

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

  const handleCreateTimeline = async () => {
    try {
      console.log('Creating timeline entry with data:', newTimelineData);

      // Convert date to ISO format
      const isoDate = new Date(newTimelineData.date).toISOString();

      // Create timeline entry
      await api.createTimelineEntry({
        name: newTimelineData.title,
        text: newTimelineData.description,
        location: newTimelineData.location,
        date: isoDate,
      });

      // Refresh timeline entries
      await fetchTimelineEntries();

      // Reset form
      setNewTimelineData({ date: "", title: "", description: "", location: "" });
      setShowCreateTimeline(false);
    } catch (err) {
      console.error('Failed to create timeline entry:', err);
      alert('Falha ao criar entrada. Verifique o console para mais detalhes.');
    }
  };

  const handleDeleteEntry = async () => {
    if (!entryToDelete) return;

    try {
      setIsDeletingEntry(true);
      setDeleteError(null);

      await api.deleteTimelineEntry(entryToDelete.id);

      // Close modal
      setEntryToDelete(null);

      // Refresh the timeline
      await fetchTimelineEntries();
    } catch (err) {
      console.error('Failed to delete timeline entry:', err);
      setDeleteError('Falha ao excluir a entrada. Tente novamente.');
    } finally {
      setIsDeletingEntry(false);
    }
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

          {loading ? (
            <div className="text-center py-16">
              <p className="text-muted-foreground">Carregando...</p>
            </div>
          ) : error ? (
            <div className="text-center py-16">
              <p className="text-destructive mb-4">{error}</p>
              <Button onClick={fetchTimelineEntries}>Tentar Novamente</Button>
            </div>
          ) : timelineEvents.length === 0 ? (
            <div className="text-center py-16 bg-muted rounded-lg">
              <Calendar className="h-16 w-16 mx-auto text-muted-foreground/50 mb-4" />
              <h3 className="text-lg font-medium text-foreground mb-2">Nenhuma entrada na linha do tempo ainda</h3>
              <p className="text-muted-foreground mb-4">
                Adicione a primeira entrada para começar a contar a história!
              </p>
              <ShowWhenAuthenticated>
                <Button onClick={() => setShowCreateTimeline(true)}>
                  Adicionar Primeira Entrada
                </Button>
              </ShowWhenAuthenticated>
            </div>
          ) : (
            <div className="relative">
              {/* Continuous timeline line */}
              <div
                className="absolute left-[30px] top-[30px] w-1 bg-primary"
                style={{ height: `calc(100% - 60px)` }}
              ></div>

              {timelineEvents.map((event) => (
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
                        <div className="flex items-center gap-2 shrink-0">
                          <Badge variant="outline" className="flex items-center gap-1">
                            <MapPin className="h-3 w-3" />
                            {event.location}
                          </Badge>
                          <ShowWhenAuthenticated>
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-8 w-8 text-destructive hover:text-destructive hover:bg-destructive/10"
                              onClick={() => setEntryToDelete(event)}
                            >
                              <X className="h-4 w-4" />
                            </Button>
                          </ShowWhenAuthenticated>
                        </div>
                      </CardTitle>
                    </CardHeader>
                    <CardContent>
                      <p className="text-muted-foreground leading-relaxed">
                        {event.description}
                      </p>
                    </CardContent>
                  </Card>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Create Timeline Entry Modal */}
        <CreateTimelineEntryModal
          open={showCreateTimeline}
          onOpenChange={setShowCreateTimeline}
          timelineData={newTimelineData}
          onTimelineDataChange={setNewTimelineData}
          onSubmit={handleCreateTimeline}
        />

        {/* Delete Timeline Entry Confirmation Modal */}
        {entryToDelete && (
          <Dialog open={!!entryToDelete} onOpenChange={() => setEntryToDelete(null)}>
            <DialogContent className="sm:max-w-md">
              <DialogHeader>
                <DialogTitle>Excluir Entrada</DialogTitle>
              </DialogHeader>

              {deleteError && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{deleteError}</AlertDescription>
                </Alert>
              )}

              <div className="space-y-4">
                <p className="text-muted-foreground">
                  Tem certeza que deseja excluir a entrada <strong>{entryToDelete.title}</strong>?
                </p>
                <p className="text-sm text-muted-foreground">
                  Esta ação não pode ser desfeita. A entrada será permanentemente removida da linha do tempo.
                </p>
              </div>

              <DialogFooter className="gap-2 sm:gap-0">
                <Button
                  variant="outline"
                  onClick={() => setEntryToDelete(null)}
                  disabled={isDeletingEntry}
                >
                  Cancelar
                </Button>
                <Button
                  variant="destructive"
                  onClick={handleDeleteEntry}
                  disabled={isDeletingEntry}
                  className="gap-2"
                >
                  <Trash2 className="h-4 w-4" />
                  {isDeletingEntry ? 'Excluindo...' : 'Excluir Entrada'}
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        )}
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
  onSubmit: () => void;
}

const CreateTimelineEntryModal: React.FC<CreateTimelineEntryModalProps> = ({
  open,
  onOpenChange,
  timelineData,
  onTimelineDataChange,
  onSubmit,
}) => {

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit();
  };

  const handleClose = () => {
    onOpenChange(false);
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
    </>
  );
};

export default Historia;
