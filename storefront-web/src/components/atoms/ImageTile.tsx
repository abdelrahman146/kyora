import { memo } from 'react';
import { Image } from 'lucide-react';

interface ImageTileProps {
  src?: string;
  alt?: string;
  aspectClassName?: string;
  className?: string;
}

/**
 * ImageTile Atom - Display image with fallback
 * Memoized to prevent unnecessary re-renders
 * Uses lazy loading for performance optimization
 */
export const ImageTile = memo<ImageTileProps>(function ImageTile({
  src,
  alt = '',
  aspectClassName = 'aspect-square',
  className = '',
}) {
  return (
    <div
      className={`w-full overflow-hidden rounded-box bg-base-200 ${aspectClassName} ${className}`}
    >
      {src ? (
        <img
          className="h-full w-full object-cover"
          src={src}
          alt={alt}
          loading="lazy"
        />
      ) : (
        <div className="h-full w-full flex items-center justify-center">
          <Image className="h-8 w-8 opacity-50" strokeWidth={1.5} />
        </div>
      )}
    </div>
  );
});
