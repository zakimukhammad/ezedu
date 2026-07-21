import { useEffect } from 'preact/hooks';

interface Props {
  active?: boolean;
}

export default function SensorySparkle({ active = true }: Props) {
  useEffect(() => {
    if (!active) return;

    const createSparkle = (x: number, y: number) => {
      const colors = ['#f59e0b', '#ec4899', '#10b981', '#6366f1', '#fbbf24', '#34d399'];
      const count = 6 + Math.floor(Math.random() * 6);

      for (let i = 0; i < count; i++) {
        const particle = document.createElement('div');
        particle.className = 'sparkle-particle';
        
        const color = colors[Math.floor(Math.random() * colors.length)];
        const size = Math.floor(Math.random() * 14) + 10;
        const angle = (Math.PI * 2 * i) / count;
        const distance = Math.random() * 50 + 20;

        const dx = Math.cos(angle) * distance;
        const dy = Math.sin(angle) * distance;

        particle.style.cssText = `
          position: fixed;
          left: ${x}px;
          top: ${y}px;
          width: ${size}px;
          height: ${size}px;
          background: ${color};
          border-radius: 50%;
          pointer-events: none;
          z-index: 9999;
          box-shadow: 0 0 10px ${color};
          transition: transform 0.6s cubic-bezier(0.1, 0.8, 0.3, 1), opacity 0.6s ease-out;
          transform: translate(-50%, -50%) scale(1);
          opacity: 1;
        `;

        document.body.appendChild(particle);

        requestAnimationFrame(() => {
          particle.style.transform = `translate(calc(-50% + ${dx}px), calc(-50% + ${dy}px)) scale(0)`;
          particle.style.opacity = '0';
        });

        setTimeout(() => {
          particle.remove();
        }, 650);
      }
    };

    const handlePointerDown = (e: PointerEvent) => {
      createSparkle(e.clientX, e.clientY);
    };

    window.addEventListener('pointerdown', handlePointerDown);
    return () => {
      window.removeEventListener('pointerdown', handlePointerDown);
    };
  }, [active]);

  return null;
}
