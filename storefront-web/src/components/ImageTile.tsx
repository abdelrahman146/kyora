import { PhotoIcon } from '@heroicons/react/24/outline';

export function ImageTile(props: {
  src?: string;
  alt?: string;
  aspectClassName?: string;
  className?: string;
}) {
  const aspectClassName = props.aspectClassName || 'aspect-square';

  return (
    <div
      className={`w-full overflow-hidden rounded-box bg-base-200 ${aspectClassName} ${
        props.className || ''
      }`}
    >
      {props.src ? (
        <img
          className="h-full w-full object-cover"
          src={props.src}
          alt={props.alt || ''}
          loading="lazy"
        />
      ) : (
        <div className="h-full w-full flex items-center justify-center">
          <PhotoIcon className="h-8 w-8 opacity-50" />
        </div>
      )}
    </div>
  );
}
