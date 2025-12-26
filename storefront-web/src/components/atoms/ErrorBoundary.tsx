import React from 'react';

type Props = {
  children: React.ReactNode;
};

type State = {
  hasError: boolean;
};

/**
 * ErrorBoundary prevents the storefront UI from crashing completely on render-time errors.
 */
export class ErrorBoundary extends React.Component<Props, State> {
  public state: State = { hasError: false };

  public static getDerivedStateFromError(): State {
    return { hasError: true };
  }

  public render() {
    if (!this.state.hasError) return this.props.children;

    return (
      <div className="min-h-dvh flex items-center justify-center p-6">
        <div className="card card-border w-full max-w-lg">
          <div className="card-body">
            <h1 className="card-title">Something went wrong</h1>
            <p className="opacity-80">Please reload the page and try again.</p>
            <div className="card-actions justify-end">
              <button type="button" className="btn" onClick={() => window.location.reload()}>
                Reload
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
