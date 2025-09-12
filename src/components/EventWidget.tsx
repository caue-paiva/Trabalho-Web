import { useState, useEffect } from "react";
import { Calendar, MapPin, Users, Clock } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useLanguage } from "@/hooks/useLanguage";

interface Event {
  id: string;
  title: string;
  date: string;
  location: string;
  description: string;
  attendees?: number;
}

const EventWidget = () => {
  const [events, setEvents] = useState<Event[]>([]);
  const [loading, setLoading] = useState(true);
  const { t } = useLanguage();

  useEffect(() => {
    // Simulating API call to eventos.grupysanca.com.br
    const fetchEvents = async () => {
      try {
        // For now, using mock data since we can't access the real API
        const mockEvents: Event[] = [
          {
            id: "1",
            title: "Python para Iniciantes: Primeiros Passos",
            date: "2024-01-25T19:00:00",
            location: "IFSP São Carlos",
            description: "Uma introdução prática ao Python para quem nunca programou antes.",
            attendees: 35
          },
          {
            id: "2", 
            title: "Workshop: Análise de Dados com Pandas",
            date: "2024-02-15T18:30:00",
            location: "USP São Carlos",
            description: "Aprenda a manipular e analisar dados usando a biblioteca Pandas.",
            attendees: 28
          },
          {
            id: "3",
            title: "Encontro Mensal: Machine Learning",
            date: "2024-03-08T19:00:00", 
            location: "Coworking Central",
            description: "Discussão sobre projetos de ML e networking entre pythonistas.",
            attendees: 42
          }
        ];
        
        setEvents(mockEvents);
      } catch (error) {
        console.error("Erro ao carregar eventos:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchEvents();
  }, []);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return {
      day: date.getDate().toString().padStart(2, '0'),
      month: date.toLocaleDateString('pt-BR', { month: 'short' }),
      time: date.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
    };
  };

  if (loading) {
    return (
      <Card className="w-full">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Calendar className="h-5 w-5 text-primary" />
            {t("events.title")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="animate-pulse space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-20 bg-muted rounded-lg"></div>
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="w-full">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Calendar className="h-5 w-5 text-primary" />
          {t("events.title")}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {events.length === 0 ? (
          <p className="text-muted-foreground text-center py-8">
            {t("events.noEvents")}
          </p>
        ) : (
          events.map((event) => {
            const { day, month, time } = formatDate(event.date);
            
            return (
              <div
                key={event.id}
                className="flex gap-4 p-4 border border-border rounded-lg hover:bg-secondary/30 transition-colors"
              >
                <div className="flex flex-col items-center bg-primary text-primary-foreground rounded-lg p-3 min-w-[60px]">
                  <span className="text-xs font-medium uppercase">{month}</span>
                  <span className="text-xl font-bold">{day}</span>
                </div>
                
                <div className="flex-1 space-y-2">
                  <h3 className="font-semibold text-foreground leading-tight">
                    {event.title}
                  </h3>
                  
                  <div className="flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
                    <div className="flex items-center gap-1">
                      <Clock className="h-3 w-3" />
                      {time}
                    </div>
                    <div className="flex items-center gap-1">
                      <MapPin className="h-3 w-3" />
                      {event.location}
                    </div>
                    {event.attendees && (
                      <div className="flex items-center gap-1">
                        <Users className="h-3 w-3" />
                        {event.attendees} {t("events.attendees")}
                      </div>
                    )}
                  </div>
                  
                  <p className="text-sm text-muted-foreground line-clamp-2">
                    {event.description}
                  </p>
                  
                  <div className="flex items-center justify-between pt-2">
                    <Badge variant="secondary">{t("events.free")}</Badge>
                    <Button size="sm" variant="outline">
                      {t("events.learnMore")}
                    </Button>
                  </div>
                </div>
              </div>
            );
          })
        )}
        
        <div className="pt-4 border-t border-border">
          <Button variant="ghost" className="w-full" asChild>
            <a href="https://eventos.grupysanca.com.br" target="_blank" rel="noopener noreferrer">
              {t("events.viewAll")}
            </a>
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

export default EventWidget;