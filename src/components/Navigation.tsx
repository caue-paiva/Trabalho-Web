import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import { Menu, X, Github, Instagram } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/ThemeToggle";
import { LanguageToggle } from "@/components/LanguageToggle";
import { useLanguage } from "@/hooks/useLanguage";
import { useTheme } from "@/hooks/useTheme";
import { TelegramIcon } from "@/components/Icon";

const Navigation = () => {
  const [isOpen, setIsOpen] = useState(false);
  const location = useLocation();
  const { t } = useLanguage();
  const { actualTheme } = useTheme();

  const navItems = [
    { href: "/", label: t("navigation.home") },
    { href: "/historia", label: t("navigation.history") },
    { href: "/galeria", label: t("navigation.gallery") },
    { href: "/codigo-conduta", label: t("navigation.codeOfConduct") },
    { href: "/contato", label: t("navigation.contact") },
  ];

  const socialLinks = [
    { href: "https://github.com/grupy-sanca", icon: Github, label: "GitHub", isCustom: false },
    { href: "https://instagram.com/grupysanca", icon: Instagram, label: "Instagram", isCustom: false },
    { href: "https://t.me/grupysanca", icon: TelegramIcon, label: "Telegram", isCustom: true },
  ];

  return (
    <nav className="bg-card border-b border-border sticky top-0 z-50 backdrop-blur-sm bg-card/80">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <Link to="/" className="flex items-center space-x-3 hover:opacity-80 transition-opacity">
            <img
              src={actualTheme === "light" ? "/grupy-logo.png" : "/logo_grupy_branca.svg"}
              alt="Grupy Sanca"
              className="h-8 w-auto object-contain"
              style={{ maxWidth: '120px', height: '32px' }}
            />
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-8">
            {navItems.map((item) => (
              <Link
                key={item.href}
                to={item.href}
                className={`text-sm font-medium transition-colors hover:text-primary ${location.pathname === item.href
                  ? "text-primary border-b-2 border-primary pb-1"
                  : "text-foreground"
                  }`}
              >
                {item.label}
              </Link>
            ))}

            <div className="flex items-center space-x-2 border-l border-border pl-4">
              <ThemeToggle />
              <LanguageToggle />
              {socialLinks.map((social) => (
                <Button
                  key={social.label}
                  variant="ghost"
                  size="sm"
                  asChild
                  className="hover:bg-secondary"
                >
                  <a href={social.href} target="_blank" rel="noopener noreferrer">
                    {social.isCustom ? (
                      <TelegramIcon size={16} className="text-foreground" />
                    ) : (
                      <social.icon className="h-4 w-4" />
                    )}
                    <span className="sr-only">{social.label}</span>
                  </a>
                </Button>
              ))}
            </div>
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsOpen(!isOpen)}
              className="hover:bg-secondary"
            >
              {isOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
            </Button>
          </div>
        </div>

        {/* Mobile Navigation */}
        {isOpen && (
          <div className="md:hidden py-4 border-t border-border bg-card">
            <div className="flex flex-col space-y-3">
              {navItems.map((item) => (
                <Link
                  key={item.href}
                  to={item.href}
                  onClick={() => setIsOpen(false)}
                  className={`text-sm font-medium transition-colors hover:text-primary px-2 py-1 rounded ${location.pathname === item.href
                    ? "text-primary bg-secondary"
                    : "text-foreground"
                    }`}
                >
                  {item.label}
                </Link>
              ))}

              <div className="flex items-center space-x-2 pt-2 border-t border-border">
                <ThemeToggle />
                <LanguageToggle />
                {socialLinks.map((social) => (
                  <Button
                    key={social.label}
                    variant="ghost"
                    size="sm"
                    asChild
                    className="hover:bg-secondary"
                  >
                    <a href={social.href} target="_blank" rel="noopener noreferrer">
                      {social.isCustom ? (
                        <TelegramIcon size={16} className="text-foreground" />
                      ) : (
                        <social.icon className="h-4 w-4" />
                      )}
                      <span className="sr-only">{social.label}</span>
                    </a>
                  </Button>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navigation;