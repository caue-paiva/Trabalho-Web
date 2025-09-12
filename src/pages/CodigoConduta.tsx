import { ExternalLink, Users, Heart, Shield, AlertTriangle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";

const CodigoConduta = () => {
  const principles = [
    {
      icon: Users,
      title: "Inclusão e Diversidade",
      description: "Valorizamos a participação de todas as pessoas, independentemente de sua experiência, gênero, orientação sexual, deficiência, aparência física, raça, etnia, idade, religião ou nacionalidade."
    },
    {
      icon: Heart,
      title: "Respeito Mútuo",
      description: "Tratamos todos os participantes com dignidade e respeito, criando um ambiente onde todos se sintam seguros para aprender e compartilhar."
    },
    {
      icon: Shield,
      title: "Ambiente Seguro",
      description: "Não toleramos assédio, discriminação ou comportamento inadequado em qualquer forma, seja presencial ou online."
    }
  ];

  const expectations = [
    "Use linguagem acolhedora e inclusiva",
    "Seja respeitoso com diferentes pontos de vista e experiências",
    "Aceite críticas construtivas com elegância",
    "Foque no que é melhor para a comunidade",
    "Mostre empatia com outros membros da comunidade",
    "Ajude iniciantes e compartilhe conhecimento de forma positiva"
  ];

  const unacceptableBehaviors = [
    "Comentários ofensivos relacionados a gênero, orientação sexual, raça, religião ou deficiência",
    "Ameaças ou linguagem violenta dirigida a outras pessoas",
    "Assédio público ou privado",
    "Publicação de informações privadas de terceiros sem permissão",
    "Conduta que poderia ser considerada inadequada em um ambiente profissional",
    "Intimidação deliberada ou perseguição"
  ];

  return (
    <div className="min-h-screen bg-background py-12">
      <div className="max-w-4xl mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-6">
            Código de Conduta
          </h1>
          <p className="text-xl text-muted-foreground max-w-3xl mx-auto leading-relaxed">
            Nosso compromisso em manter uma comunidade acolhedora, diversa e respeitosa para todos.
          </p>
        </div>

        {/* Alert */}
        <Alert className="mb-8">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            Este código de conduta aplica-se a todos os eventos, espaços online e comunicações 
            relacionadas ao Grupy São Carlos.
          </AlertDescription>
        </Alert>

        {/* Introduction */}
        <div className="gradient-section p-8 rounded-xl mb-12">
          <h2 className="text-3xl font-bold text-foreground mb-6">Nossa missão</h2>
          <div className="prose prose-lg text-muted-foreground leading-relaxed">
            <p className="mb-4">
              O Grupy São Carlos tem o compromisso de ser uma comunidade livre de assédio para 
              todos, independentemente do sexo, identidade de gênero e expressão, idade, orientação 
              sexual, deficiência, aparência física, tamanho corporal, cor de pele, etnia, religião 
              (ou falta dela) ou escolhas de tecnologias.
            </p>
            <p>
              Não toleramos o assédio aos participantes da comunidade, sob qualquer forma. 
              Linguagem e imagens sexuais não são adequadas em nenhum local da comunidade, 
              incluindo palestras, workshops, festas, Twitter e outras mídias on-line.
            </p>
          </div>
        </div>

        {/* Core Principles */}
        <div className="mb-12">
          <h2 className="text-3xl font-bold text-foreground mb-8 text-center">
            Nossos valores
          </h2>
          <div className="grid md:grid-cols-3 gap-6">
            {principles.map((principle, index) => (
              <Card key={index} className="text-center hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="mx-auto bg-primary/10 p-3 rounded-full w-fit mb-4">
                    <principle.icon className="h-8 w-8 text-primary" />
                  </div>
                  <CardTitle className="text-lg">{principle.title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-muted-foreground text-sm leading-relaxed">
                    {principle.description}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>

        {/* Expected Behavior */}
        <div className="mb-12">
          <Card>
            <CardHeader>
              <CardTitle className="text-2xl text-primary">Comportamentos esperados</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="space-y-3">
                {expectations.map((expectation, index) => (
                  <li key={index} className="flex items-start gap-3">
                    <div className="w-2 h-2 bg-primary rounded-full mt-2 shrink-0"></div>
                    <span className="text-muted-foreground">{expectation}</span>
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        </div>

        {/* Unacceptable Behavior */}
        <div className="mb-12">
          <Card className="border-destructive/20">
            <CardHeader>
              <CardTitle className="text-2xl text-destructive">Comportamentos inaceitáveis</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="space-y-3">
                {unacceptableBehaviors.map((behavior, index) => (
                  <li key={index} className="flex items-start gap-3">
                    <div className="w-2 h-2 bg-destructive rounded-full mt-2 shrink-0"></div>
                    <span className="text-muted-foreground">{behavior}</span>
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        </div>

        {/* Reporting */}
        <div className="mb-12">
          <Card className="gradient-section">
            <CardHeader>
              <CardTitle className="text-2xl">Como reportar violações</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <p className="text-muted-foreground leading-relaxed">
                Se você testemunhar ou sofrer qualquer comportamento inaceitável, ou tiver outras 
                preocupações relacionadas à comunidade, entre em contato conosco imediatamente:
              </p>
              
              <div className="bg-card p-4 rounded-lg border border-border">
                <h4 className="font-semibold text-foreground mb-2">Canais de contato:</h4>
                <ul className="space-y-2 text-muted-foreground">
                  <li>• Email: contato@grupysanca.com.br</li>
                  <li>• Durante eventos: procure qualquer organizador</li>
                  <li>• Online: entre em contato via direct message nas redes sociais</li>
                </ul>
              </div>
              
              <p className="text-muted-foreground text-sm">
                Todas as denúncias serão tratadas com confidencialidade e seriedade. 
                Tomaremos as medidas apropriadas para lidar com a situação de forma justa e transparente.
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Link to GitHub */}
        <div className="text-center">
          <Card className="inline-block">
            <CardContent className="p-6">
              <h3 className="text-xl font-semibold text-foreground mb-4">
                Versão completa e oficial
              </h3>
              <p className="text-muted-foreground mb-6">
                Para a versão mais atualizada e detalhada do nosso Código de Conduta, 
                consulte nosso repositório oficial no GitHub.
              </p>
              <Button size="lg" asChild>
                <a 
                  href="https://github.com/grupy-sanca/codigo-de-conduta" 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="gap-2"
                >
                  Ver no GitHub
                  <ExternalLink className="h-4 w-4" />
                </a>
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default CodigoConduta;