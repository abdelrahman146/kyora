import {
  useEffect,
  useMemo,
  useState,
  type MouseEventHandler,
  type ReactNode,
} from "react";

export interface ResendCountdownButtonProps {
  cooldownSeconds: number;
  isBusy?: boolean;
  onResend: () => void | Promise<void>;
  renderLabel: (args: { remainingSeconds: number; canResend: boolean }) => ReactNode;
  className?: string;
}

export function ResendCountdownButton({
  cooldownSeconds,
  isBusy = false,
  onResend,
  renderLabel,
  className,
}: ResendCountdownButtonProps) {
  const normalizedCooldown = useMemo(() => {
    if (!Number.isFinite(cooldownSeconds)) return 0;
    return Math.max(0, Math.floor(cooldownSeconds));
  }, [cooldownSeconds]);

  const [remainingSeconds, setRemainingSeconds] = useState(normalizedCooldown);

  // Reset countdown whenever the cooldown value changes.
  useEffect(() => {
    setRemainingSeconds(normalizedCooldown);
  }, [normalizedCooldown]);

  const canResend = remainingSeconds <= 0;

  useEffect(() => {
    if (canResend) return;

    const timer = setInterval(() => {
      setRemainingSeconds((s) => Math.max(0, s - 1));
    }, 1000);

    return () => {
      clearInterval(timer);
    };
  }, [canResend]);

  const handleClick: MouseEventHandler<HTMLButtonElement> = (e) => {
    e.preventDefault();
    if (!canResend || isBusy) return;
    void onResend();
  };

  return (
    <button
      type="button"
      className={className ?? "btn btn-ghost btn-sm"}
      disabled={!canResend || isBusy}
      onClick={handleClick}
    >
      {renderLabel({ remainingSeconds, canResend })}
    </button>
  );
}
