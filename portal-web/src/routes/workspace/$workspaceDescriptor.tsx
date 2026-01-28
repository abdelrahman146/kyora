import { Outlet, createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/workspace/$workspaceDescriptor')({
  component: WorkspaceLayout,
})

function WorkspaceLayout() {
  return <Outlet />
}
