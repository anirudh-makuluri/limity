import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Limity Next.js Example',
  description: 'Rate limiting example with Limity',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
