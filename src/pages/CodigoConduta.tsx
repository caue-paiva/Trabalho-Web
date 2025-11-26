import { ExternalLink, AlertTriangle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription } from "@/components/ui/alert";

const CodigoConduta = () => {

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
            relacionadas ao grupy-sanca.
          </AlertDescription>
        </Alert>

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