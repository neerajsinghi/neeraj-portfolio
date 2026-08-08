export type AnalyticsParams = Record<string, string | number | boolean>;

declare global {
  interface Window {
    dataLayer?: IArguments[];
    gtag?: (command: "event", eventName: string, params?: AnalyticsParams) => void;
  }
}

export function trackEvent(eventName: string, params: AnalyticsParams = {}) {
  if (typeof window === "undefined") return;
  if (typeof window.gtag !== "function") {
    window.dataLayer = window.dataLayer || [];
    window.gtag = function gtag() {
      window.dataLayer?.push(arguments);
    };
  }
  window.gtag("event", eventName, params);
}