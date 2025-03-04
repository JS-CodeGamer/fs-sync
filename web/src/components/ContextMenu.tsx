import React from 'react';
import { Link, Download, Copy, Trash2 } from 'lucide-react';

interface ContextMenuProps {
  x: number;
  y: number;
  onClose: () => void;
  onAction: (action: string) => void;
}

export default function ContextMenu({ x, y, onClose, onAction }: ContextMenuProps) {
  const menuItems = [
    { icon: Link, label: 'Copy link', action: 'copyLink' },
    { icon: Download, label: 'Download', action: 'download' },
    { icon: Copy, label: 'Make a copy', action: 'copy' },
    { icon: Trash2, label: 'Move to trash', action: 'trash' },
  ];

  return (
    <>
      <div
        className="fixed inset-0"
        onClick={onClose}
      />
      <div
        className="fixed bg-white rounded-lg shadow-lg py-2 w-56"
        style={{ top: y, left: x }}
      >
        {menuItems.map(({ icon: Icon, label, action }) => (
          <button
            key={action}
            onClick={() => {
              onAction(action);
              onClose();
            }}
            className="w-full px-4 py-2 text-left flex items-center space-x-2 hover:bg-gray-100"
          >
            <Icon className="h-4 w-4 text-gray-500" />
            <span>{label}</span>
          </button>
        ))}
      </div>
    </>
  );
}
