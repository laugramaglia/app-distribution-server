import {
  BarChart,
  BookOpen,
  Calendar,
  ClipboardCheck,
  Clock,
  Code,
  HeartPulse,
  MessageSquare,
  Plane,
  Receipt,
  Shield,
  icons,
} from 'lucide-react';
import type { IconName } from './types';

export const iconMap = {
  Clock,
  Receipt,
  ClipboardCheck,
  Calendar,
  Shield,
  Code,
  BarChart,
  BookOpen,
  HeartPulse,
  Plane,
  MessageSquare,
};

export function getIconComponent(iconId: IconName) {
  return iconMap[iconId as keyof typeof iconMap] || icons['Package'];
}
