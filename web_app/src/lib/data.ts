import type { App } from './types';
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
} from 'lucide-react';

const apps: App[] = [
  {
    id: 'timetrack',
    name: 'TimeTrack',
    description: 'Track your work hours and project time efficiently.',
    category: 'Productivity',
    client: 'Internal',
    icon: Clock,
    iconId: 'Clock',
    versions: [
      {
        version: '2.1.0',
        releaseDate: '2024-05-20',
        changelog: 'Added project-based time tracking and reporting features.',
        downloadUrl: '#',
      },
      {
        version: '2.0.5',
        releaseDate: '2024-03-15',
        changelog: 'Bug fixes and performance improvements.',
        downloadUrl: '#',
      },
    ],
  },
  {
    id: 'expenseit',
    name: 'ExpenseIt',
    description: 'Submit and manage your expense reports with ease.',
    category: 'Finance',
    client: 'Internal',
    icon: Receipt,
    iconId: 'Receipt',
    versions: [
      {
        version: '1.5.2',
        releaseDate: '2024-06-01',
        changelog: 'Improved receipt scanning accuracy.',
        downloadUrl: '#',
      },
    ],
  },
  {
    id: 'taskmaster',
    name: 'TaskMaster',
    description: 'Organize personal and team tasks, set deadlines, and track progress.',
    category: 'Productivity',
    client: 'Client A',
    icon: ClipboardCheck,
    iconId: 'ClipboardCheck',
    versions: [
      {
        version: '3.0.0',
        releaseDate: '2024-06-10',
        changelog: 'Major UI overhaul with new collaboration features.',
        downloadUrl: '#',
      },
      {
        version: '2.9.8',
        releaseDate: '2024-05-02',
        changelog: 'Fixed sync issues on mobile devices.',
        downloadUrl: '#',
      },
    ],
  },
  {
    id: 'meetingsched',
    name: 'MeetingSched',
    description: 'Schedule meetings, check availability, and book rooms.',
    category: 'Collaboration',
    client: 'Internal',
    icon: Calendar,
    iconId: 'Calendar',
    versions: [
      {
        version: '1.8.0',
        releaseDate: '2024-04-25',
        changelog: 'Integration with third-party calendar services.',
        downloadUrl: '#',
      },
    ],
  },
  {
    id: 'securepass',
    name: 'SecurePass',
    description: 'Securely store and manage your passwords and sensitive information.',
    category: 'Security',
    client: 'Client B',
    icon: Shield,
    iconId: 'Shield',
    versions: [
      {
        version: '4.2.1',
        releaseDate: '2024-06-12',
        changelog: 'Enhanced encryption and added support for security keys.',
        downloadUrl: '#',
      },
    ],
  },
  {
    id: 'codeshare',
    name: 'CodeShare',
    description: 'Share code snippets and collaborate on projects with your team.',
    category: 'Development',
    client: 'Internal',
    icon: Code,
    iconId: 'Code',
    versions: [
      {
        version: '1.2.0',
        releaseDate: '2024-05-18',
        changelog: 'Added syntax highlighting for 10 new languages.',
        downloadUrl: '#',
      },
    ],
  },
  {
    id: 'dataview',
    name: 'DataView',
    description: 'Visualize and analyze company data with interactive dashboards.',
    category: 'Analytics',
    client: 'Client A',
    icon: BarChart,
    iconId: 'BarChart',
    versions: [
      {
        version: '2.5.0',
        releaseDate: '2024-06-05',
        changelog: 'New chart types and data export options.',
        downloadUrl: '#',
      },
    ],
  },
  {
    id: 'learnit',
    name: 'LearnIt',
    description: 'Access learning resources, courses, and training materials.',
    category: 'HR & Training',
    client: 'Internal',
    icon: BookOpen,
    iconId: 'BookOpen',
    versions: [
      {
        version: '1.1.3',
        releaseDate: '2024-03-30',
        changelog: 'Updated course catalog and improved video player.',
        downloadUrl: '#',
      },
    ],
  },
  {
    id: 'healthplus',
    name: 'HealthPlus',
    description: 'Access health and wellbeing resources and company benefits.',
    category: 'Wellness',
    client: 'Internal',
    icon: HeartPulse,
    iconId: 'HeartPulse',
    versions: [
      {
        version: '1.0.0',
        releaseDate: '2024-02-01',
        changelog: 'Initial release.',
        downloadUrl: '#',
      },
    ],
  },
  {
    id: 'travelbuddy',
    name: 'TravelBuddy',
    description: 'Plan and book work-related travel and accommodations.',
    category: 'Travel',
    client: 'Client C',
    icon: Plane,
    iconId: 'Plane',
    versions: [
      {
        version: '2.3.0',
        releaseDate: '2024-05-22',
        changelog: 'Added flight tracking and loyalty program integration.',
        downloadUrl: '#',
      },
    ],
  },
  {
    id: 'feedbackloop',
    name: 'FeedbackLoop',
    description: 'Share and collect feedback with colleagues and managers.',
    category: 'Collaboration',
    client: 'Internal',
    icon: MessageSquare,
    iconId: 'MessageSquare',
    versions: [
      {
        version: '1.4.0',
        releaseDate: '2024-06-15',
        changelog: 'Anonymous feedback option and 360-degree reviews.',
        downloadUrl: '#',
      },
    ],
  },
];

export async function getApps(): Promise<App[]> {
  // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 200));
  return apps;
}

export async function getAppById(id: string): Promise<App | undefined> {
    // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 100));
  return apps.find(app => app.id === id);
}

export async function getAppByName(name: string): Promise<App | undefined> {
  await new Promise(resolve => setTimeout(resolve, 100));
  return apps.find(app => app.name.toLowerCase() === name.toLowerCase());
}
