import { SiTelegram } from "react-icons/si";

export function TelegramIcon({
    size = 20,
    className = "inline-block",
}: { size?: number; className?: string }) {
    return <SiTelegram size={size} className={className} aria-hidden="true" />;
}

