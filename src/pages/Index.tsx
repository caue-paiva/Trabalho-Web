import { ArrowRight, Code, Users, Calendar, Heart, Clock, MapPin } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useLanguage } from "@/hooks/useLanguage";
import { useTheme } from "@/hooks/useTheme";
import { useState, useEffect } from "react";

interface Event {
  id: string;
  title: string;
  date: string;
  location: string;
  description: string;
  attendees?: number;
}

const Index = () => {
  const { t } = useLanguage();
  const { actualTheme } = useTheme();
  const [featuredEvent, setFeaturedEvent] = useState<Event | null>(null);

  useEffect(() => {
    // Fetch a single upcoming event
    const fetchFeaturedEvent = async () => {
      try {
        // Using mock data similar to EventWidget
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

        // Get the first upcoming event
        setFeaturedEvent(mockEvents[0]);
      } catch (error) {
        console.error("Erro ao carregar evento:", error);
      }
    };

    fetchFeaturedEvent();
  }, []);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return {
      day: date.getDate().toString().padStart(2, '0'),
      month: date.toLocaleDateString('pt-BR', { month: 'short' }),
      time: date.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
    };
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <section className="gradient-hero py-20 px-4">
        <div className="max-w-6xl mx-auto">
          <div className="grid lg:grid-cols-2 gap-8 items-end">
            {/* Left side - Logo, Title, Subtitle, Buttons */}
            <div className="text-center lg:text-left">
              <div className="float mb-8">
                <img
                  src={actualTheme === "light" ? "/grupy-logo.png" : "/logo_grupy_branca.svg"}
                  alt="Grupy Sanca"
                  className="h-24 w-auto mx-auto lg:mx-0 mb-6 object-contain"
                  style={{ maxWidth: '360px', height: '96px' }}
                />
              </div>

              <h1 className="text-4xl md:text-6xl font-bold text-foreground mb-6 fade-in">
                {t("home.title")}
              </h1>

              <p className="text-xl md:text-2xl text-muted-foreground mb-8 leading-relaxed">
                {t("home.subtitle")}
              </p>

              <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
                <Button size="lg" className="px-8 py-3 text-lg" asChild>
                  <a href="/historia">
                    {t("home.learnHistory")}
                    <ArrowRight className="ml-2 h-5 w-5" />
                  </a>
                </Button>

                <Button size="lg" className="px-8 py-3 text-lg text-black font-semibold hover:opacity-90" style={{ backgroundColor: '#FADA5E' }} asChild>
                  <a href="https://eventos.grupysanca.com.br" target="_blank" rel="noopener noreferrer">
                    {t("home.viewEvents")}
                    <Calendar className="ml-2 h-5 w-5" />
                  </a>
                </Button>
              </div>
            </div>

            {/* Right side - Featured Event */}
            <div className="hidden lg:flex lg:flex-col lg:items-end lg:justify-end">
              {featuredEvent ? (() => {
                const { day, month, time } = formatDate(featuredEvent.date);
                return (
                  <div className="w-3/5 flex flex-col gap-2">
                    {/* Header box */}
                    <div className="bg-secondary/50 rounded-lg px-4 py-2">
                      <h3 className="text-sm font-semibold text-foreground">Próximo Evento</h3>
                    </div>

                    {/* Event card */}
                    <div className="min-h-[280px]">
                      <div className="flex flex-col gap-4 p-6 border border-border rounded-lg hover:bg-secondary/30 transition-colors h-full bg-card">
                        <div className="flex gap-4">
                          <div className="flex flex-col items-center bg-primary text-primary-foreground rounded-lg p-3 min-w-[60px]">
                            <span className="text-xs font-medium uppercase">{month}</span>
                            <span className="text-xl font-bold">{day}</span>
                          </div>

                          <div className="flex-1 space-y-3">
                            <h3 className="font-semibold text-foreground leading-tight">
                              {featuredEvent.title}
                            </h3>

                            <div className="flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
                              <div className="flex items-center gap-1">
                                <Clock className="h-3 w-3" />
                                {time}
                              </div>
                              <div className="flex items-center gap-1">
                                <MapPin className="h-3 w-3" />
                                {featuredEvent.location}
                              </div>
                              {featuredEvent.attendees && (
                                <div className="flex items-center gap-1">
                                  <Users className="h-3 w-3" />
                                  {featuredEvent.attendees} {t("events.attendees")}
                                </div>
                              )}
                            </div>

                            <p className="text-sm text-muted-foreground line-clamp-3">
                              {featuredEvent.description}
                            </p>
                          </div>
                        </div>

                        <div className="flex items-center justify-between pt-2 mt-auto">
                          <Badge variant="secondary">{t("events.free")}</Badge>
                          <Button size="sm" variant="outline" asChild>
                            <a href="https://eventos.grupysanca.com.br" target="_blank" rel="noopener noreferrer">
                              {t("events.learnMore")}
                            </a>
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })() : (
                <div className="w-3/5 flex flex-col gap-2">
                  <div className="bg-secondary/50 rounded-lg px-4 py-2">
                    <h3 className="text-sm font-semibold text-foreground">Próximo Evento</h3>
                  </div>
                  <div className="min-h-[280px]">
                    <div className="animate-pulse space-y-4 p-6">
                      <div className="h-full bg-muted rounded-lg"></div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </section>

      {/* About Section */}
      <section id="sobre" className="py-20 px-4 gradient-section">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-6">
              {t("home.whoWeAre")}
            </h2>
            <p className="text-lg text-muted-foreground max-w-3xl mx-auto leading-relaxed">
              {t("home.description")}
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8 mb-16">
            <Card className="text-center hover:shadow-lg transition-shadow">
              <CardHeader>
                <div className="mx-auto bg-primary/10 p-3 rounded-full w-fit mb-4">
                  <Users className="h-8 w-8 text-primary" />
                </div>
                <CardTitle>{t("home.community.title")}</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  {t("home.community.description")}
                </p>
              </CardContent>
            </Card>

            <Card className="text-center hover:shadow-lg transition-shadow">
              <CardHeader>
                <div className="mx-auto bg-secondary/10 p-3 rounded-full w-fit mb-4">
                  <Code className="h-8 w-8 text-secondary" />
                </div>
                <CardTitle>{t("home.learning.title")}</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  {t("home.learning.description")}
                </p>
              </CardContent>
            </Card>

            <Card className="text-center hover:shadow-lg transition-shadow">
              <CardHeader>
                <div className="mx-auto bg-accent/10 p-3 rounded-full w-fit mb-4">
                  <Heart className="h-8 w-8 text-accent" />
                </div>
                <CardTitle>{t("home.networking.title")}</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  {t("home.networking.description")}
                </p>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 px-4 gradient-hero">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-6">
            {t("home.joinCommunity")}
          </h2>
          <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
            {t("home.joinDescription")}
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button size="lg" className="px-8 py-3 text-lg" asChild>
              <a href="/historia">
                {t("home.learnHistory")}
                <ArrowRight className="ml-2 h-5 w-5" />
              </a>
            </Button>

            <Button variant="outline" size="lg" className="px-8 py-3 text-lg" asChild>
              <a href="/galeria">
                {t("home.seePhotos")}
              </a>
            </Button>
          </div>
        </div>
      </section>
    </div>
  );
};

export default Index;