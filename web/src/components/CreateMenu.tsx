import React from 'react';
import { FolderPlus, Upload, File } from 'lucide-react';

interface CreateMenuProps {
  x: number;
  y: number;
  onClose: () => void;
  onAction: (action: string) => void;
}

export default function CreateMenu({ x, y, onClose, onAction }: CreateMenuProps) {
  const menuItems = [
    { icon: FolderPlus, label: 'New Folder', action: 'newFolder' },
    { icon: Upload, label: 'Upload Folder', action: 'uploadFolder' },
    { icon: File, label: 'Upload File', action: 'uploadFile' },
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