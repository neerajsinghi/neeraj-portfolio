export type Repo = {
  name: string;
  description: string;
  html_url: string;
  language: string;
  stargazers_count: number;
  forks_count: number;
};

export type ChatItem =
  | { kind: "user"; text: string }
  | { kind: "bot"; text: string; running?: boolean }
  | { kind: "trace"; tool: string; arg: string; sources?: string[] }
  | { kind: "typing" };
