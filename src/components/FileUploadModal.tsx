/**
 * FileUploadModal Component
 *
 * A reusable modal component for uploading files with drag-and-drop support.
 * Can be configured to accept different file types (images, documents, videos, etc.).
 */

import { useState } from "react";
import { Upload, X, FileIcon, ImageIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";

export interface FileUploadConfig {
  /** Accepted file types (e.g., "image/*", "application/pdf", ".doc,.docx") */
  accept?: string;
  /** Maximum file size in bytes (default: 10MB) */
  maxSize?: number;
  /** Maximum number of files allowed (default: unlimited) */
  maxFiles?: number;
  /** Allow multiple file selection (default: true) */
  multiple?: boolean;
  /** Custom validation function */
  validate?: (file: File) => { valid: boolean; error?: string };
  /** File type category for display purposes */
  fileCategory?: "image" | "document" | "video" | "audio" | "any";
}

interface FileUploadModalProps {
  /** Controls modal visibility */
  open: boolean;
  /** Callback when modal visibility changes */
  onOpenChange: (open: boolean) => void;
  /** Callback when files are uploaded */
  onUpload: (files: File[]) => void;
  /** Modal title */
  title?: string;
  /** Upload button text */
  uploadButtonText?: string;
  /** Configuration options */
  config?: FileUploadConfig;
}

interface FileWithError {
  file: File;
  error?: string;
}

const DEFAULT_CONFIG: FileUploadConfig = {
  accept: "*/*",
  maxSize: 10 * 1024 * 1024, // 10MB
  multiple: true,
  fileCategory: "any",
};

export const FileUploadModal: React.FC<FileUploadModalProps> = ({
  open,
  onOpenChange,
  onUpload,
  title = "Upload Files",
  uploadButtonText = "Upload",
  config = {},
}) => {
  const [dragActive, setDragActive] = useState(false);
  const [selectedFiles, setSelectedFiles] = useState<FileWithError[]>([]);

  const finalConfig = { ...DEFAULT_CONFIG, ...config };

  const validateFile = (file: File): { valid: boolean; error?: string } => {
    // Check file size
    if (finalConfig.maxSize && file.size > finalConfig.maxSize) {
      return {
        valid: false,
        error: `File too large (max ${(finalConfig.maxSize / 1024 / 1024).toFixed(1)}MB)`,
      };
    }

    // Check file type if accept is specified
    if (finalConfig.accept && finalConfig.accept !== "*/*") {
      const acceptTypes = finalConfig.accept.split(",").map((t) => t.trim());
      const fileType = file.type;
      const fileName = file.name;

      const isAccepted = acceptTypes.some((acceptType) => {
        if (acceptType.startsWith(".")) {
          // Extension-based check
          return fileName.toLowerCase().endsWith(acceptType.toLowerCase());
        } else if (acceptType.includes("/*")) {
          // Wildcard type check (e.g., "image/*")
          const baseType = acceptType.split("/")[0];
          return fileType.startsWith(baseType + "/");
        } else {
          // Exact type check
          return fileType === acceptType;
        }
      });

      if (!isAccepted) {
        return { valid: false, error: "File type not accepted" };
      }
    }

    // Custom validation
    if (finalConfig.validate) {
      return finalConfig.validate(file);
    }

    return { valid: true };
  };

  const addFiles = (files: File[]) => {
    const filesToAdd: FileWithError[] = [];

    for (const file of files) {
      // Check max files limit
      if (
        finalConfig.maxFiles &&
        selectedFiles.length + filesToAdd.length >= finalConfig.maxFiles
      ) {
        break;
      }

      const validation = validateFile(file);
      filesToAdd.push({
        file,
        error: validation.valid ? undefined : validation.error,
      });
    }

    setSelectedFiles((prev) => [...prev, ...filesToAdd]);
  };

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    const files = Array.from(e.dataTransfer.files);
    addFiles(files);
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const files = Array.from(e.target.files);
      addFiles(files);
    }
  };

  const removeFile = (index: number) => {
    setSelectedFiles((prev) => prev.filter((_, i) => i !== index));
  };

  const handleUpload = () => {
    const validFiles = selectedFiles
      .filter((f) => !f.error)
      .map((f) => f.file);

    if (validFiles.length > 0) {
      onUpload(validFiles);
      setSelectedFiles([]);
      onOpenChange(false);
    }
  };

  const handleClose = () => {
    setSelectedFiles([]);
    onOpenChange(false);
  };

  const getFileIcon = () => {
    switch (finalConfig.fileCategory) {
      case "image":
        return ImageIcon;
      default:
        return FileIcon;
    }
  };

  const getAcceptDescription = () => {
    const maxSizeMB = finalConfig.maxSize
      ? (finalConfig.maxSize / 1024 / 1024).toFixed(1)
      : "10";

    switch (finalConfig.fileCategory) {
      case "image":
        return `PNG, JPG, GIF até ${maxSizeMB}MB`;
      case "document":
        return `PDF, DOC, DOCX até ${maxSizeMB}MB`;
      case "video":
        return `MP4, MOV, AVI até ${maxSizeMB}MB`;
      case "audio":
        return `MP3, WAV, OGG até ${maxSizeMB}MB`;
      default:
        return `Arquivos até ${maxSizeMB}MB`;
    }
  };

  const IconComponent = getFileIcon();
  const validFilesCount = selectedFiles.filter((f) => !f.error).length;
  const hasErrors = selectedFiles.some((f) => f.error);

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Upload className="h-5 w-5" />
            {title}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          {/* Drag and Drop Area */}
          <div
            className={`relative border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
              dragActive
                ? "border-primary bg-primary/5"
                : "border-muted-foreground/25 hover:border-muted-foreground/50"
            }`}
            onDragEnter={handleDrag}
            onDragLeave={handleDrag}
            onDragOver={handleDrag}
            onDrop={handleDrop}
          >
            <input
              type="file"
              id="file-upload"
              multiple={finalConfig.multiple}
              accept={finalConfig.accept}
              onChange={handleFileSelect}
              className="hidden"
            />

            <div className="flex flex-col items-center gap-2">
              <IconComponent className="h-12 w-12 text-muted-foreground/50" />
              <div>
                <p className="text-sm font-medium text-foreground">
                  Arraste arquivos aqui ou{" "}
                  <label
                    htmlFor="file-upload"
                    className="text-primary hover:underline cursor-pointer"
                  >
                    escolha arquivos
                  </label>
                </p>
                <p className="text-xs text-muted-foreground mt-1">
                  {getAcceptDescription()}
                </p>
                {finalConfig.maxFiles && (
                  <p className="text-xs text-muted-foreground">
                    Máximo de {finalConfig.maxFiles} arquivo(s)
                  </p>
                )}
              </div>
            </div>
          </div>

          {/* Selected Files List */}
          {selectedFiles.length > 0 && (
            <div className="space-y-2">
              <h4 className="text-sm font-medium text-foreground">
                Arquivos selecionados ({selectedFiles.length})
                {hasErrors && (
                  <span className="text-destructive ml-2">
                    ({selectedFiles.filter((f) => f.error).length} com erro)
                  </span>
                )}
              </h4>
              <div className="max-h-48 overflow-y-auto space-y-2">
                {selectedFiles.map((fileWithError, index) => (
                  <div
                    key={index}
                    className={`flex items-center justify-between p-2 rounded-lg ${
                      fileWithError.error ? "bg-destructive/10" : "bg-muted"
                    }`}
                  >
                    <div className="flex items-center gap-2 flex-1 min-w-0">
                      <IconComponent
                        className={`h-4 w-4 flex-shrink-0 ${
                          fileWithError.error
                            ? "text-destructive"
                            : "text-muted-foreground"
                        }`}
                      />
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <span
                            className={`text-sm truncate ${
                              fileWithError.error
                                ? "text-destructive"
                                : "text-foreground"
                            }`}
                          >
                            {fileWithError.file.name}
                          </span>
                          <span className="text-xs text-muted-foreground flex-shrink-0">
                            ({(fileWithError.file.size / 1024).toFixed(1)} KB)
                          </span>
                        </div>
                        {fileWithError.error && (
                          <p className="text-xs text-destructive">
                            {fileWithError.error}
                          </p>
                        )}
                      </div>
                    </div>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => removeFile(index)}
                      className="flex-shrink-0"
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex gap-2 justify-end pt-4">
            <Button variant="outline" onClick={handleClose}>
              Cancelar
            </Button>
            <Button
              onClick={handleUpload}
              disabled={validFilesCount === 0}
              className="gap-2"
            >
              <Upload className="h-4 w-4" />
              {uploadButtonText}{" "}
              {validFilesCount > 0 && `(${validFilesCount})`}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};
