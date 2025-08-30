import type { ComponentType } from 'react';
import type { icons } from 'lucide-react';

export type Version = {
  version: string;
  releaseDate: string;
  changelog: string;
  downloadUrl: string;
};

export type IconName = keyof typeof icons;

export type App = {
  id: string;
  name: string;
  description: string;
  category: string;
  client: string;
  icon: ComponentType<{ className?: string }>;
  iconId: IconName; // Add this
  versions: Version[];
};

export type SerializableApp = Omit<App, 'icon'>;
