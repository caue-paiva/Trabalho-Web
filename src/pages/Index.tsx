import { ArrowRight, Code, Users, Calendar, Heart } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import EventWidget from "@/components/EventWidget";
import { useLanguage } from "@/hooks/useLanguage";
import grupyLogo from "@/assets/grupy-logo.png";

const Index = () => {
  const { t } = useLanguage();
  
  return (
    <div className="min-h-screen bg-background">
      {/* Hero Section */}
      <section className="gradient-hero py-20 px-4">
        <div className="max-w-6xl mx-auto text-center">
          <div className="float mb-8">
            <img 
              src={grupyLogo} 
              alt="Grupy São Carlos" 
              className="h-24 w-auto mx-auto mb-6"
            />
          </div>
          
          <h1 className="text-4xl md:text-6xl font-bold text-foreground mb-6 fade-in">
            {t("home.title")}
          </h1>
          
          <p className="text-xl md:text-2xl text-muted-foreground mb-8 max-w-3xl mx-auto leading-relaxed">
            {t("home.subtitle")}
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <Button size="lg" className="px-8 py-3 text-lg" asChild>
              <a href="#sobre">
                {t("home.learnHistory")}
                <ArrowRight className="ml-2 h-5 w-5" />
              </a>
            </Button>
            
            <Button variant="outline" size="lg" className="px-8 py-3 text-lg" asChild>
              <a href="https://eventos.grupysanca.com.br" target="_blank" rel="noopener noreferrer">
                {t("home.viewEvents")}
                <Calendar className="ml-2 h-5 w-5" />
              </a>
            </Button>
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

      {/* Events Section */}
      <section className="py-20 px-4">
        <div className="max-w-6xl mx-auto">
          <div className="grid lg:grid-cols-2 gap-12 items-start">
            <div>
              <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-6">
                {t("home.upcomingEvents")}
              </h2>
              <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
                {t("home.eventsDescription")}
              </p>
              
              <div className="space-y-4 mb-8">
                <div className="flex items-center gap-3">
                  <div className="w-2 h-2 bg-primary rounded-full"></div>
                  <span className="text-muted-foreground">
                    <strong>{t("home.features.workshops")}</strong>
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  <div className="w-2 h-2 bg-secondary rounded-full"></div>
                  <span className="text-muted-foreground">
                    <strong>{t("home.features.talks")}</strong>
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  <div className="w-2 h-2 bg-accent rounded-full"></div>
                  <span className="text-muted-foreground">
                    <strong>{t("home.features.networking")}</strong>
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  <div className="w-2 h-2 bg-primary rounded-full"></div>
                  <span className="text-muted-foreground">
                    <strong>{t("home.features.dojos")}</strong>
                  </span>
                </div>
              </div>
            </div>

            <div>
              <EventWidget />
            </div>
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