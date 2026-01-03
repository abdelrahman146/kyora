/**
 * Atoms - Basic UI Building Blocks
 *
 * Atomic components are the smallest, most reusable UI elements.
 * They have minimal or no dependencies on other components.
 */

// Layout & Display
export { ErrorBoundary } from './ErrorBoundary'
export { Button, type ButtonProps } from './Button'
export { Badge, type BadgeProps } from './Badge'
export {
  Skeleton,
  SkeletonText,
  type SkeletonProps,
  type SkeletonTextProps,
} from './Skeleton'
export { Avatar, type AvatarProps } from './Avatar'
export { IconButton, type IconButtonProps } from './IconButton'
export { Logo, type LogoProps } from './Logo'
export { Dropdown, type DropdownProps } from './Dropdown'
export { Modal, type ModalProps } from './Modal'
export { Dialog } from './Dialog'
export { SocialMediaLink, type SocialMediaLinkProps } from './SocialMediaLink'
export { type SocialPlatform } from '../icons/social'

// Form Components (TanStack Form integrated)
export { FormInput, type FormInputProps } from './FormInput'
export { FormTextarea, type FormTextareaProps } from './FormTextarea'
export { FormCheckbox, type FormCheckboxProps } from './FormCheckbox'
export { FormToggle, type FormToggleProps } from './FormToggle'
export {
  FormRadio,
  type FormRadioProps,
  type FormRadioOption,
} from './FormRadio'
export {
  FormSelect,
  type FormSelectProps,
  type FormSelectOption,
} from './FormSelect'
export { PriceInput, type PriceInputProps } from './PriceInput'
export { PasswordInput, type PasswordInputProps } from './PasswordInput'
export { DatePicker, type DatePickerProps } from './DatePicker'
export { TimePicker, type TimePickerProps } from './TimePicker'
export { DateRangePicker, type DateRangePickerProps } from './DateRangePicker'
export { DateTimePicker, type DateTimePickerProps } from './DateTimePicker'

// File Upload Components
export { FormFileInput, type FormFileInputProps } from './FormFileInput'
export { FileUploadZone, type FileUploadZoneProps } from './FileUploadZone'
export { FilePreview, type FilePreviewProps } from './FilePreview'
export { UploadProgress, type UploadProgressProps } from './UploadProgress'
export { FileThumbnail, type FileThumbnailProps } from './FileThumbnail'
export { DragHandle, type DragHandleProps } from './DragHandle'
