import React from "react";
import { Folder, File, MoreVertical } from "lucide-react";

interface FileItem {
  id: string;
  name: string;
  type: "folder" | "file";
  size?: string;
  modifiedAt: string;
}

interface FileListProps {
  items: FileItem[] | null;
  onContextMenu: (e: React.MouseEvent, item: FileItem) => void;
}

export default function FileList({ items, onContextMenu }: FileListProps) {
  return (
    <div className="w-full">
      <div className="grid grid-cols-[auto_1fr_1fr_auto] gap-4 px-6 py-2 text-sm text-gray-500 border-b">
        <div className="w-8"></div>
        <div>Name</div>
        <div>Last modified</div>
        <div>File size</div>
      </div>

      {items?.map((item) => (
        <div
          key={item.id}
          onContextMenu={(e) => onContextMenu(e, item)}
          className="grid grid-cols-[auto_1fr_1fr_auto] gap-4 px-6 py-2 hover:bg-gray-100 cursor-pointer items-center"
        >
          <div className="w-8">
            {item.type === "folder" ? (
              <Folder className="h-5 w-5 text-gray-400" />
            ) : (
              <File className="h-5 w-5 text-gray-400" />
            )}
          </div>
          <div>{item.name}</div>
          <div>{item.modifiedAt}</div>
          <div className="flex items-center space-x-2">
            <span>{item.size}</span>
            <button className="p-1 hover:bg-gray-200 rounded-full">
              <MoreVertical className="h-4 w-4 text-gray-500" />
            </button>
          </div>
        </div>
      ))}
    </div>
  );
}
