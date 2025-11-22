import { ArrowRight, Code, Users, Calendar, Heart, Clock, MapPin } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useLanguage } from "@/hooks/useLanguage";
import { useTheme } from "@/hooks/useTheme";
import { useState, useEffect } from "react";
import { ExternalEvent, getLatestExternalEvent } from "@/services/api";
import DOMPurify from "dompurify";

const Index = () => {
  const { t } = useLanguage();
  const { actualTheme } = useTheme();
  const [featuredEvent, setFeaturedEvent] = useState<ExternalEvent | null>(null);

  useEffect(() => {
    const loadFeaturedEvent = async () => {
      try {
        const ev = await getLatestExternalEvent();
        setFeaturedEvent(ev);
      } catch (err) {
        console.error("Failed to load latest external event", err);
      }
    };

    loadFeaturedEvent();
  }, []);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return {
      day: date.getDate().toString().padStart(2, '0'),
      month: date.toLocaleDateString('pt-BR', { month: 'short' }),
      time: date.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
    };
  };

  const sanitizeAndExtractText = (html: string | undefined): string => {
    if (!html) return '';
    // Convert HTML line breaks to newlines before stripping tags
    const text = html
      .replace(/<br\s*\/?>/gi, '\n')           // <br>, <br/>, <br />
      .replace(/<\/p>/gi, '\n\n')              // </p> to double newline
      .replace(/<p[^>]*>/gi, '');              // Remove opening <p> tags

    // Sanitize to remove any remaining HTML tags and dangerous content
    const clean = DOMPurify.sanitize(text, { ALLOWED_TAGS: [] });

    // Return trimmed text with preserved newlines
    return clean.trim();
  };

  const linkifyText = (text: string) => {
    // URL regex pattern that matches http(s) URLs and t.me links
    const urlRegex = /(https?:\/\/[^\s]+|t\.me\/[^\s)]+)/g;
    const parts = text.split(urlRegex);

    return parts.map((part, index) => {
      // Check if this part is a URL
      if (part.match(urlRegex)) {
        // Add https:// protocol if it's a t.me link without protocol
        const href = part.startsWith('t.me/') ? `https://${part}` : part;
        return (
          <a
            key={index}
            href={href}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 dark:text-blue-400 hover:underline"
          >
            {part}
          </a>
        );
      }
      // Return regular text, preserving newlines
      return <span key={index}>{part}</span>;
    });
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
                const { day, month, time } = formatDate(featuredEvent.starts_at);
                return (
                  <div className="w-[90%] flex flex-col gap-2">
                    {/* Header box */}
                    <div className="bg-secondary/50 rounded-lg px-4 py-2">
                      <h3 className="text-sm font-semibold text-foreground">Próximo Evento</h3>
                    </div>

                    {/* Event card */}
                    <div>
                      <div className="flex flex-col gap-4 p-6 border border-border rounded-lg hover:bg-secondary/30 transition-colors bg-card">
                        <div className="flex gap-4">
                          <div className="flex flex-col items-center bg-primary text-primary-foreground rounded-lg p-3 min-w-[60px]">
                            <span className="text-xs font-medium uppercase">{month}</span>
                            <span className="text-xl font-bold">{day}</span>
                          </div>

                          <div className="flex-1 min-w-0 space-y-3">
                            <h3 className="font-semibold text-foreground leading-tight break-words">
                              {featuredEvent.name}
                            </h3>

                            <div className="flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
                              {/* clock, location, etc */}
                              <div className="flex items-center gap-1">
                                <Clock className="h-3 w-3" />
                                {time}
                              </div>
                              {featuredEvent.location_name && (
                                <div className="flex items-center gap-1">
                                  <MapPin className="h-3 w-3" />
                                  {featuredEvent.location_name}
                                </div>
                              )}
                            </div>

                            {/* removed line-clamp-3 so the box can grow with content */}
                            <p className="text-sm text-muted-foreground break-words whitespace-pre-line">
                              {linkifyText(sanitizeAndExtractText(featuredEvent.description))}
                            </p>
                          </div>
                        </div>

                        <div className="flex items-center justify-between pt-2 mt-auto">
                          <Badge variant="secondary">{t("events.free")}</Badge>
                          <Button size="sm" variant="outline" asChild>
                            <a
                              href={featuredEvent.link ?? "https://eventos.grupysanca.com.br"}
                              target="_blank"
                              rel="noopener noreferrer"
                            >
                              {t("events.learnMore")}
                            </a>
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>
                );
              })() : (
                <div className="w-[90%] flex flex-col gap-2">
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