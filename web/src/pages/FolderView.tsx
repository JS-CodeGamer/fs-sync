import React, { useState, useRef, useEffect } from "react";
import { useNavigate, useParams } from "react-router";

import Navbar from "../components/Navbar";
import FileList from "../components/FileList";
import ContextMenu from "../components/ContextMenu";
import CreateMenu from "../components/CreateMenu";
import FolderDialog from "../components/FolderDialog";
import { Asset } from "../models/asset";
import backend from "../helpers/backend";
import { useAuth } from "../context/AuthProvider";
import { toast } from "react-toastify";

type Props = { assetID: string };

export function FolderView({ assetID }: Props) {
  const navigate = useNavigate();

  const [assetInfo, setAssetInfo] = useState<Asset>({} as Asset);

  const { isLoggedIn } = useAuth();

  const [contextMenu, setContextMenu] = useState<{
    x: number;
    y: number;
    item: Asset | null;
  } | null>(null);
  const [createMenu, setCreateMenu] = useState<{ x: number; y: number } | null>(
    null,
  );
  const [showFolderDialog, setShowFolderDialog] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const folderInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    backend
      .get<Asset>(`/asset/${assetID}`)
      .then((res) => {
        if (!res.data) navigate("/404");
        setAssetInfo(res.data);
      })
      .catch(() => {
        toast.error("Unable to fetch data!!!");
      });
  }, []);

  const handleContextMenu = (e: React.MouseEvent, item: Asset) => {
    e.preventDefault();
    setContextMenu({ x: e.clientX, y: e.clientY, item });
    setCreateMenu(null);
  };

  const handleBackgroundContextMenu = (e: React.MouseEvent) => {
    if ((e.target as HTMLElement).classList.contains("background-area")) {
      e.preventDefault();
      setCreateMenu({ x: e.clientX, y: e.clientY });
      setContextMenu(null);
    }
  };

  const handleContextMenuAction = (action: string) => {
    if (!contextMenu?.item) return;

    switch (action) {
      case "copyLink":
        console.log("Copy link for:", contextMenu.item.name);
        break;
      case "download":
        console.log("Download:", contextMenu.item.name);
        break;
      case "copy":
        const newItem = {
          ...contextMenu.item,
          id: Date.now().toString(),
          name: `${contextMenu.item.name} (Copy)`,
        };
        setAssetInfo((assetInfo) => ({
          ...assetInfo,
          children: [...(assetInfo?.children ?? []), newItem],
        }));
        break;
      case "trash":
        setAssetInfo(
          assetInfo.filter((item) => item.id !== contextMenu.item?.id),
        );
        break;
    }
    setContextMenu(null);
  };

  const handleCreateMenuAction = (action: string) => {
    switch (action) {
      case "newFolder":
        setShowFolderDialog(true);
        break;
      case "uploadFile":
        fileInputRef.current?.click();
        break;
      case "uploadFolder":
        folderInputRef.current?.click();
        break;
    }
    setCreateMenu(null);
  };

  const handleCreateFolder = (name: string) => {
    const newFolder = {
      name,
      parent: assetID,
      is_dir: true,
    };
    backend
      .post("/a", newFolder)
      .then((res) => {
        if (res.data && res.data.message === "success")
          setAssetInfo((x) => ({
            ...x,
            children: [
              ...x?.children,
              {
                ...newFolder,
              },
            ],
          }));
        else throw Error("error");
      })
      .catch(() => {
        toast.error("Unable to create folder!!!");
      })
      .finally(() => {
        setShowFolderDialog(false);
      });
  };

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);
    const newItems = files.map((file) => ({
      id: Date.now().toString() + Math.random(),
      name: file.name,
      type: "file" as const,
      size: `${(file.size / (1024 * 1024)).toFixed(1)} MB`,
      modifiedAt: new Date().toLocaleDateString("en-US", {
        month: "short",
        day: "numeric",
        year: "numeric",
      }),
    }));
    setAssetInfo([...assetInfo, ...newItems]);
    e.target.value = "";
  };

  if (!isLoggedIn()) {
    navigate("/login?return=" + encodeURIComponent(window.location.pathname));
  }

  return (
    <div className="min-h-screen bg-white">
      <Navbar />

      <div
        className="background-area p-6 min-h-[calc(100vh-64px)]"
        onContextMenu={handleBackgroundContextMenu}
      >
        <FileList
          parent={assetInfo?.children}
          onContextMenu={handleContextMenu}
        />
      </div>

      {contextMenu && (
        <ContextMenu
          x={contextMenu.x}
          y={contextMenu.y}
          onClose={() => setContextMenu(null)}
          onAction={handleContextMenuAction}
        />
      )}

      {createMenu && (
        <CreateMenu
          x={createMenu.x}
          y={createMenu.y}
          onClose={() => setCreateMenu(null)}
          onAction={handleCreateMenuAction}
        />
      )}

      {showFolderDialog && (
        <FolderDialog
          onClose={() => setShowFolderDialog(false)}
          onSubmit={handleCreateFolder}
        />
      )}

      <input
        type="file"
        ref={fileInputRef}
        className="hidden"
        onChange={handleFileUpload}
        multiple
      />
      <input
        type="file"
        ref={folderInputRef}
        className="hidden"
        onChange={handleFileUpload}
        /* @ts-expect-error */
        webkitdirectory=""
        multiple
      />
    </div>
  );
}

export function RouteBasedFolderView() {
  const navigate = useNavigate();
  const { assetID } = useParams();
  if (!assetID) navigate("/404");
  return <FolderView assetID={assetID!} />;
}

export function RootAssetFolderView() {
  const navigate = useNavigate();
  const { isLoggedIn, user } = useAuth();
  if (!isLoggedIn() || !user?.root_asset) navigate("/404");
  return <FolderView assetID={user?.root_asset!} />;
}
