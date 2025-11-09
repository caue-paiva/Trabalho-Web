import { Calendar, MapPin, Users, Award } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

const Historia = () => {
  const timelineEvents = [
    {
      year: "2015",
      title: "Fundação do Grupy Sanca",
      description: "O grupo nasceu da necessidade de criar uma comunidade local de Python em São Carlos, inspirado por outros grupys pelo Brasil.",
      participants: 12,
      location: "IFSP São Carlos"
    },
    {
      year: "2016",
      title: "Primeiro Workshop Oficial",
      description: "Organizamos nosso primeiro workshop sobre Django, marcando o início das atividades educacionais regulares.",
      participants: 25,
      location: "USP São Carlos"
    },
    {
      year: "2017",
      title: "Parcerias com Universidades",
      description: "Estabelecemos parcerias formais com USP e IFSP para sediar eventos e atingir mais estudantes.",
      participants: 45,
      location: "Múltiplas instituições"
    },
    {
      year: "2018",
      title: "Primeiro Python Day São Carlos",
      description: "Organizamos um evento de dia inteiro com palestras, workshops e networking, nosso maior evento até então.",
      participants: 120,
      location: "Centro de Convenções"
    },
    {
      year: "2019",
      title: "Expansão Regional",
      description: "O grupo começou a receber participantes de cidades vizinhas, consolidando-se como referência regional.",
      participants: 80,
      location: "São Carlos e região"
    },
    {
      year: "2020-2021",
      title: "Eventos Online",
      description: "Durante a pandemia, adaptamos todos os eventos para formato online, mantendo a comunidade ativa e unida.",
      participants: 150,
      location: "Online"
    },
    {
      year: "2022",
      title: "Retorno Híbrido",
      description: "Retomamos os eventos presenciais com transmissão online, ampliando nosso alcance e inclusividade.",
      participants: 200,
      location: "Híbrido"
    },
    {
      year: "2023",
      title: "8 Anos de Comunidade",
      description: "Celebramos 8 anos de atividades contínuas, com mais de 50 eventos realizados e centenas de pessoas impactadas.",
      participants: 250,
      location: "São Carlos"
    }
  ];

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
          <h2 className="text-3xl font-bold text-foreground mb-8 text-center">
            Linha do tempo
          </h2>
          <div className="space-y-8">
            {timelineEvents.map((event, index) => (
              <div key={index} className="flex gap-6 items-start">
                <div className="flex flex-col items-center">
                  <div className="bg-primary text-primary-foreground rounded-full p-3 font-bold text-lg min-w-[60px] text-center">
                    {event.year}
                  </div>
                  {index < timelineEvents.length - 1 && (
                    <div className="w-px bg-border h-24 mt-4"></div>
                  )}
                </div>

                <Card className="flex-1 hover:shadow-lg transition-shadow">
                  <CardHeader>
                    <CardTitle className="flex items-start justify-between flex-wrap gap-2">
                      <span>{event.title}</span>
                      <div className="flex gap-2">
                        <Badge variant="secondary" className="flex items-center gap-1">
                          <Users className="h-3 w-3" />
                          {event.participants}
                        </Badge>
                        <Badge variant="outline" className="flex items-center gap-1">
                          <MapPin className="h-3 w-3" />
                          {event.location}
                        </Badge>
                      </div>
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <p className="text-muted-foreground leading-relaxed mb-4">
                      {event.description}
                    </p>
                    {/* Image placeholders - exclude first entry (2015) and 2020-2021 entry */}
                    {index !== 0 && event.year !== "2020-2021" && (
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

        {/* Photos Section */}
        <div className="gradient-section p-8 rounded-xl">
          <h2 className="text-3xl font-bold text-foreground mb-6 text-center">
            Momentos históricos
          </h2>
          <p className="text-center text-muted-foreground mb-8">
            Algumas fotos marcantes da nossa trajetória ao longo dos anos.
          </p>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[1, 2, 3, 4, 5, 6].map((index) => (
              <div key={index} className="bg-muted rounded-lg h-48 flex items-center justify-center">
                <div className="text-center text-muted-foreground">
                  <Calendar className="h-12 w-12 mx-auto mb-2 opacity-50" />
                  <p className="text-sm">Foto histórica #{index}</p>
                  <p className="text-xs">Em breve disponível</p>
                </div>
              </div>
            ))}
          </div>

          <div className="text-center mt-8">
            <p className="text-muted-foreground">
              Quer contribuir com fotos históricas? Entre em contato conosco!
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Historia;