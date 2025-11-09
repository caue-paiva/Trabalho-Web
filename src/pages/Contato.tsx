import { Github, Mail, Instagram, ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useLanguage } from "@/hooks/useLanguage";
import { TelegramIcon } from "@/components/Icon";

const Contato = () => {
    const { t } = useLanguage();

    const contactLinks = [
        {
            icon: Github,
            label: t("contact.github.label"),
            href: "https://github.com/grupy-sanca",
            description: t("contact.github.description"),
            color: "text-foreground"
        },
        {
            icon: Mail,
            label: t("contact.email.label"),
            href: "mailto:contato@grupysanca.com.br",
            description: t("contact.email.description"),
            color: "text-blue-500"
        },
        {
            icon: Instagram,
            label: t("contact.instagram.label"),
            href: "https://instagram.com/grupysanca",
            description: t("contact.instagram.description"),
            color: "text-pink-500"
        },
        {
            icon: TelegramIcon,
            label: t("contact.telegram.label"),
            href: "https://t.me/grupysanca",
            description: t("contact.telegram.description"),
            color: "text-blue-400"
        }
    ];

    return (
        <div className="min-h-screen bg-background py-12">
            <div className="max-w-4xl mx-auto px-4">
                {/* Header */}
                <div className="text-center mb-12">
                    <h1 className="text-4xl md:text-5xl font-bold text-foreground mb-6">
                        {t("contact.title")}
                    </h1>
                    <p className="text-xl text-muted-foreground max-w-3xl mx-auto leading-relaxed">
                        {t("contact.subtitle")}
                    </p>
                </div>

                {/* Contact Cards */}
                <div className="grid md:grid-cols-2 gap-6">
                    {contactLinks.map((contact, index) => (
                        <Card key={index} className="hover:shadow-lg transition-shadow">
                            <CardHeader>
                                <div className="flex items-center gap-4 mb-2">
                                    <div className={`p-3 rounded-lg bg-secondary/10 ${contact.color}`}>
                                        {contact.icon === TelegramIcon ? (
                                            <TelegramIcon size={24} className={contact.color} />
                                        ) : (
                                            <contact.icon className="h-6 w-6" />
                                        )}
                                    </div>
                                    <CardTitle className="text-xl">{contact.label}</CardTitle>
                                </div>
                            </CardHeader>
                            <CardContent className="space-y-4">
                                <p className="text-muted-foreground text-sm leading-relaxed">
                                    {contact.description}
                                </p>
                                <Button asChild className="w-full">
                                    <a
                                        href={contact.href}
                                        target={contact.href.startsWith("mailto:") ? undefined : "_blank"}
                                        rel={contact.href.startsWith("mailto:") ? undefined : "noopener noreferrer"}
                                        className="flex items-center justify-center gap-2"
                                    >
                                        {t("contact.visitLink")}
                                        <ExternalLink className="h-4 w-4" />
                                    </a>
                                </Button>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default Contato;

