// vendor
import React, { FC, useState, useEffect, useRef } from 'react';
import { NativeTypes } from 'react-dnd-html5-backend'
import { useDrop, useDrag } from 'react-dnd'
import File from '../file/File';
import { ListObjectsV2Command, CopyObjectCommand, DeleteObjectCommand, PutObjectCommand } from "@aws-sdk/client-s3";
import classNames from 'classnames';
// components
import { ReactComponent as FolderIcon } from 'Images/icons/folder.svg';
import { ReactComponent as FolderExpandedIcon } from 'Images/icons/folder-open.svg';
import NewFolder from '../newFolder/NewFolder';
import { InputText } from 'Components/form/text/index';
import Actions from '../actions/Actions';
// css
import './Folder.scss';

interface Folder {
  Prefix: string;
}

interface FolderContents {
  files: FileArray;
  folders: {
    [folder: string]: Folder;
  }
}

interface FileArray {
  [fileName: string]: File
}

interface FolderObject {
  [folderName: string]: {
    isSelected: boolean;
    isExpanded: boolean;
  }
}

interface FileObject {
  [fileName: string]: {
    file: File;
    isSelected: boolean;
  }
}

interface File {
  Key: string;
  Size: number;
  LastModified: any;
}

interface Props {
  folder: {
    Prefix: string;
  };
  setSelected: any;
  folderSelected: boolean;
  setExpanded: any;
  s3Client: any;
  bucket: string,
  setSingleFolder: (folderKey: string, isSelected: boolean, isExpanded: boolean) => void;
  setMultipleFolders: (folderKeys: Array<string>, isSelected: boolean, ignoreIfExists?: boolean) => void;
  setSingleFile: (fileKey: string, file: File, isSelected?: boolean) => void;
  setMultipleFiles: (fileKeys: Array<string>, files: FileArray, isSelected: boolean) => void;
  folderList: any;
  fileList: any;
  setFocusedFile: any;
  focusedFile: File,
  humanFileSize: any;
  removeFiles: (fileKeys: Array<string>) => void;
  removeFolders: (folderKeys: Array<string>) => void;
  refetchS3: any;
  folderExpanded: boolean;
  siblings: any;
  depth: number;
  setCanDropFolder: (folders: any) => void;
  canDropFolder: any;
  s3CopyFolder: (bucket: string, source: string, dest: string) => void;
  setFocusedFileData: (file: any) => void;
  batchUpload:(files: any, prefix: string) => void;
  sort: (files: any, isFolder: boolean) => any[];
  setIsLocked: (locked: boolean) => void;
  folderReader: (folder: any, fileArray: any[], path: string) => void;
  namespace: string;
  searchMode: boolean;
  setSearchFolderExpanded: any;
  effectiveFolderList: any;
  isLocked: boolean;
}

export const Folder:FC<Props> = ({
  folder,
  setExpanded,
  folderExpanded,
  setSelected,
  folderSelected,
  s3Client,
  bucket,
  setSingleFolder,
  setMultipleFolders,
  setSingleFile,
  setMultipleFiles,
  folderList,
  fileList,
  setFocusedFile,
  focusedFile,
  humanFileSize,
  removeFiles,
  removeFolders,
  refetchS3,
  siblings,
  depth,
  s3CopyFolder,
  setCanDropFolder,
  canDropFolder,
  setFocusedFileData,
  batchUpload,
  sort,
  setIsLocked,
  folderReader,
  namespace,
  searchMode,
  setSearchFolderExpanded,
  effectiveFolderList,
  isLocked,
}: Props) => {

  const splitName = folder.Prefix.split('/');
  const folderName = splitName[splitName.length - 2];
  const parentFolder = `${splitName.slice(0, -2).join('/')}/`;

  const inputRef = useRef(null);
  const checkboxRef = useRef(null);
  const nameRef = useRef(null);

  const [renameMode, setRenameMode] = useState(false);
  const [loadingState, setLoadingState] = useState(null);
  const [editedFolderName, setFolderName] = useState('');
  const [newFolderVisible, setNewFolderVisible] = useState(null);

  useEffect(() => {
    if (checkboxRef.current && !searchMode) {
      checkboxRef.current.style.marginLeft = `${(folder.Prefix.split('/').length - 3) * 30}px`;
    } else if(nameRef.current) {
      nameRef.current.style.marginLeft = `${(folder.Prefix.split('/').length - 3) * 30}px`;
    }
  }, [checkboxRef.current, nameRef.current])


  const handleRename = () => {
    setRenameMode(false);
    if (editedFolderName === folderName) {
      return;
    } else {
      s3CopyFolder(bucket, folder.Prefix, `${parentFolder}${editedFolderName}/`);
    }
  }

  const updateFolderName = (evt: any) => {
    if((evt.key === 'Enter') && (editedFolderName.length > 0)){
      handleRename();
    } else {
      setFolderName(evt.target.value);
    }
  }

  const refetchS3Local = (prefix: string) => {
    const ListObjectsInput = {
        Bucket: bucket,
        Delimiter: '/',
        Prefix: prefix,
      }
    const command = new ListObjectsV2Command(ListObjectsInput);
    s3Client.send(command)
    .then((res: any) => {
      const newFilesObject: FileArray = {};
      const fileKeys = (res.Contents && res.Contents.map((content: any) => {
        newFilesObject[content.Key] = content;
        return content.Key;
      })) || [];
      setMultipleFiles(fileKeys, newFilesObject, null);
      const folderKeys = (res.CommonPrefixes && res.CommonPrefixes.map((content: any) => content.Prefix)) || [];
      setMultipleFolders(folderKeys, null, null);
    })
  }

  const [{ isDragging }, drag, preview] = useDrag(
    () => ({
      type: 'FILE',
      item: { folder, folderName, parentFolder },
      collect: (monitor) => ({
        isDragging: !!monitor.isDragging(),
      }),
      end: (item, monitor) => {
        if (!monitor.didDrop()) {
          return
        }
      },
    }),
    [],
  )

    const [{ isOver, canDrop }, drop] = useDrop(
    () => ({
      accept: [NativeTypes.FILE, 'FILE'],
      drop: (item: any, monitor) => {
        if (item.sourcePath === folder.Prefix) {
          return;
        }
        if (item.dataTransfer && item.dataTransfer.items) {
          const { items } = item.dataTransfer;
          const fileArray: any[] = [];
          Array.from(items).forEach((item: any) =>{
            setIsLocked(true);
            const uploadedFolder = item.webkitGetAsEntry();
            if (uploadedFolder && uploadedFolder.isDirectory) {
              folderReader(uploadedFolder, fileArray, folder.Prefix)
            } else {
              fileArray.push(item.getAsFile());
            }
          })
          if (fileArray.length === items.length) {
            batchUpload(fileArray, folder.Prefix)
          }
          return;
        }
        if (item.file) {
          const copyInput = {
            Bucket: bucket,
            CopySource: `${bucket}/${item.file.Key}`,
            Key: `${folder.Prefix}${item.fileName}`,
          }
          setIsLocked(true);
          const command = new CopyObjectCommand(copyInput);
          s3Client.send(command)
            .then((res: any) => {
              const deleteCommand = new DeleteObjectCommand({
                Key: item.file.Key,
                Bucket: bucket,
              })
              s3Client.send(deleteCommand)
              .then(() => {
                refetchS3();
                refetchS3Local(folder.Prefix);
                removeFiles([item.file.Key]);
                setIsLocked(false);
              });
            });
        }
        if (item.folder
          && item.parentFolder !== folder.Prefix
          && item.folder.Prefix !== folder.Prefix) {
          s3CopyFolder(bucket, item.folder.Prefix, `${folder.Prefix}${item.folderName}/`);
        }
        if (item.files) {
          batchUpload(item.files, folder.Prefix);
        }
      },
      collect: (monitor: any) => {
        return {
          isOver: !!monitor.isOver(),
          canDrop: !!monitor.canDrop()
            &&  (monitor.getItem().sourcePath !== folder.Prefix)
            &&  (monitor.getItem().parentFolder !== folder.Prefix)
            &&  (((monitor.getItem().file )||monitor.getItem().folder && (monitor.getItem().folder.Prefix !== folder.Prefix)) || monitor.getItem().files),
        }
      },
    }),
    [],
  )

  useEffect(() => {
    setCanDropFolder((prevState: any) => ({
      ...prevState,
      [folder.Prefix]: !!(isOver && canDrop)
    }))
  }, [isOver, canDrop])

  const folderContents: FolderContents = {
    files: {},
    folders: {},
  }

  let someSelected = false;

  Object.keys(fileList).forEach((fileKey: string) => {
    if(fileKey.startsWith(folder.Prefix)) {
      folderContents.files[fileKey] = fileList[fileKey];
      if (fileList[fileKey].isSelected) {
        someSelected = true;
      }
    }
  });
  Object.keys(folderList).forEach((folderKey: string) => {
    if(folderKey.startsWith(folder.Prefix) && folderKey !== folder.Prefix) {
      folderContents.folders[folderKey] = folderList[folderKey];
      if (folderList[folderKey].isSelected) {
        someSelected = true;
      }
    }
  });

  let files = Object.keys(folderContents.files).map((fileKey) => fileList[fileKey].file).filter((file) => file.Key.split('/').length === depth);
  let folders = Object.keys(folderContents.folders).map((folderKey) => ({ Prefix: folderKey })).filter((folder) => folder.Prefix.split('/').length === (depth + 1));

  files = sort(files, false);
  folders = sort(folders, true);

  const updateState = () => {
    if(loadingState) {
      setLoadingState((prevData: any) => {
        if(prevData && !prevData.forceFinish) {
          const newCurrentProgress = prevData.currentProgress + prevData.step;
          const arcTangent = Math.atan(newCurrentProgress);
          const halfPi = Math.PI / 2;
          const progress = (Math.round((arcTangent / halfPi) * 100 * 1000) / 1000).toFixed(2);
          setTimeout(() => {
            updateState();
          }, 100);
          return {
            ...prevData,
            progress,
            currentProgress: newCurrentProgress,
          }
        } else if(prevData && prevData.forceFinish) {
          setTimeout(() => {
            setLoadingState(null);
          }, 500);
          return prevData;
        }
      })
    }
  }

  useEffect(() => {
    if(loadingState && loadingState.progress === 0) {
      updateState();
    }
  }, [loadingState])

  useEffect(() => {
    if (folderExpanded && !searchMode) {
      if (files.length === 0 && folders.length === 0) {
        setLoadingState({
          progress: 0,
          currentProgress: 0,
          step: 0.03,
          forceFinish: false,
        });
      }
      const command = new ListObjectsV2Command({
        Bucket: bucket,
        Delimiter: '/',
        Prefix: folder.Prefix,
      })
      s3Client.send(command)
      .then((res: any) => {
        const newFilesObject: FileArray = {};
        const fileKeys = (res.Contents && res.Contents.map((content: any) => {
          newFilesObject[content.Key] = content;
          return content.Key;
        })) || [];
        setMultipleFiles(fileKeys, newFilesObject, folderSelected ? folderSelected : null);
        const folderKeys = (res.CommonPrefixes && res.CommonPrefixes.map((content: any) => content.Prefix)) || [];
        setMultipleFolders(folderKeys, folderSelected ? folderSelected : null, null);
        if (fileKeys.length == 0 && folderKeys.length === 0) {
          setLoadingState((prevData: any) =>({
            ...prevData,
            progress: '100',
            forceFinish: true,
          }));
        }
      })
    }
  }, [folderExpanded])

  useEffect(() => {
    if((files.length > 0 || folders.length > 0) &&( loadingState && !loadingState.forceFinish)) {
        setLoadingState((prevData: any) =>({
          ...prevData,
          progress: '100',
          forceFinish: true,
        }));
    }
  }, [files, folders])

  useEffect(() => {
    if(renameMode) {
      inputRef.current.focus();
      inputRef.current.value = folderName;
    }
  }, [renameMode])

  useEffect(() => {
    if(newFolderVisible && !folderExpanded) {
      setExpanded();
    }
  }, [newFolderVisible])

  const checkboxCSS = classNames({
    Folder__checkbox: true,
    'Folder__checkbox--selected': folderSelected,
    'Folder__checkbox--partial': someSelected && !folderSelected,
  });
  const folderCSS = classNames({
    Folder: true,
    'Folder--expanded': folderExpanded,
    'Folder--hovered': isOver && canDrop && !isLocked,
    'Folder--dragging': isDragging,
  })

  return (
    <>
      <div
      ref={searchMode ? null : drag}
      >
        <div className={folderCSS}
          onClick={setExpanded}
          ref={searchMode ? null : drop}
        >
          {
            !searchMode && (
            <div
              ref={checkboxRef}
              className={checkboxCSS}
              onClick={() => { setSelected(folderExpanded, folderContents, siblings) }}
            />
            )
          }
          {
            searchMode && <div className="Folder__checkbox--placeholder" />
          }
          <div className="Folder__file" ref={nameRef}>
            {
              !folderExpanded && (
                <FolderIcon />
              )
            }
            {
              folderExpanded && (
                <FolderExpandedIcon />
              )
            }
            {renameMode && (
              <InputText
                inputRef={inputRef}
                label=""
                placeholder={folderName}
                updateValue={updateFolderName}
              />
            )}
            {!renameMode && `${folderName}/`}
          </div>
          <div className="Folder__size">
            {' '}
          </div>
          <div className="Folder__modified">
            {' '}
          </div>
          <div className="Folder__action">
            <Actions
              s3Client={s3Client}
              bucket={bucket}
              folder
              folderKey={folder.Prefix}
              refetchS3={() => refetchS3Local(folder.Prefix)}
              setRenameMode={() => {setRenameMode(!renameMode)}}
              renameMode={renameMode}
              handleRename={handleRename}
              removeFiles={removeFiles}
              removeFolders={removeFolders}
              allowRename={(editedFolderName.length > 0)}
              setNewFolderVisible={setNewFolderVisible}
              searchMode={searchMode}
            />
          </div>
          {
            loadingState && (
              <div className="Folder__progress-container">
                <span
                  className="Folder__progress"
                  style={{
                    width: `${loadingState.progress}%`
                  }}
                />
              </div>
            )
          }
        </div>
      </div>
      {
        folderExpanded && (
          <div className="Folder__nested">
            {
              newFolderVisible && (
                <NewFolder
                  closeFolder={() => setNewFolderVisible(false)}
                  isVisible={newFolderVisible}
                  setSingleFolder={setSingleFolder}
                  prefix={folder.Prefix}
                />
              )
            }
            {
              folders && folders.map((content: Folder) => {
              return(
                <Folder
                  key={content.Prefix}
                  folder={content}
                  isLocked={isLocked}
                  setSelected={(folderExpanded: boolean, folderContents: FolderContents, siblings: any) => {
                    if (folderList[folder.Prefix].isSelected && folderList[content.Prefix].isSelected) {
                      setSingleFolder(folder.Prefix, false, null);
                    }

                    if (!folderList[content.Prefix].isSelected && siblings) {
                      let folderSelected = true;
                      siblings.files.forEach((file: File) => {
                        if (!fileList[file.Key].isSelected) {
                          folderSelected = false;
                        }
                      })
                      siblings.folders.forEach((folder: Folder) => {
                        if(content.Prefix === folder.Prefix) {
                          return;
                        }
                        if (!folderList[folder.Prefix].isSelected) {
                          folderSelected = false;
                        }
                      })
                      if (folderSelected) {
                        setSingleFolder(folder.Prefix, true, null);
                      }
                    }

                    if (folderExpanded) {
                      const fileKeys = Object.keys(folderContents.files);
                      const files: FileArray  = {};
                      fileKeys.forEach((fileKey) => files[fileKey] = fileList[fileKey].file);
                      setMultipleFiles(fileKeys, files, !folderList[content.Prefix].isSelected);
                      const folders = Object.keys(folderContents.folders);
                      setMultipleFolders(folders, !folderList[content.Prefix].isSelected, null);
                    }
                    setSingleFolder(content.Prefix, !folderList[content.Prefix].isSelected, null)
                    }
                  }
                  folderSelected={folderList[content.Prefix].isSelected}
                  setExpanded={(evt: Event) => {
                    if(searchMode) {
                      setSearchFolderExpanded(content.Prefix, !effectiveFolderList[content.Prefix].isExpanded);
                      return;
                    }
                    if (!evt) {
                      setSingleFolder(content.Prefix, null, !folderList[content.Prefix].isExpanded);
                      return;
                    }
                    if((typeof (evt.target as HTMLElement).className) !== 'string') {
                      return;
                    }
                    if (!((evt.target as HTMLElement).className.indexOf('Folder__checkbox') > -1) && (evt.target as HTMLElement).tagName === 'DIV') {
                      setSingleFolder(content.Prefix, null, !folderList[content.Prefix].isExpanded)
                    }
                  }}
                  folderExpanded={folderList[content.Prefix].isExpanded}
                  s3Client={s3Client}
                  bucket={bucket}
                  setSingleFolder={setSingleFolder}
                  setMultipleFolders={setMultipleFolders}
                  setSingleFile={setSingleFile}
                  setMultipleFiles={setMultipleFiles}
                  folderList={folderList}
                  fileList={fileList}
                  setFocusedFile={setFocusedFile}
                  focusedFile={focusedFile}
                  humanFileSize={humanFileSize}
                  removeFiles={removeFiles}
                  removeFolders={removeFolders}
                  refetchS3={refetchS3}
                  siblings={{files, folders}}
                  depth={depth + 1}
                  setCanDropFolder={setCanDropFolder}
                  canDropFolder={canDropFolder}
                  s3CopyFolder={s3CopyFolder}
                  setFocusedFileData={setFocusedFileData}
                  batchUpload={batchUpload}
                  sort={sort}
                  setIsLocked={setIsLocked}
                  folderReader={folderReader}
                  namespace={namespace}
                  searchMode={searchMode}
                  setSearchFolderExpanded={setSearchFolderExpanded}
                  effectiveFolderList={effectiveFolderList}
                />
              )
              })
            }
            {
              files && files.map((content: File) => {
              return(
                <File
                  key={content.Key}
                  file={content}
                  setFocus={(evt: Event) => {
                    if((typeof (evt.target as HTMLElement).className) !== 'string') {
                      return;
                    }
                    if (!((evt.target as HTMLElement).className.indexOf('File__checkbox') > -1) && (evt.target as HTMLElement).tagName === 'DIV') {
                      if (focusedFile && content.Key === focusedFile.Key) {
                        setFocusedFile(null);
                        setFocusedFileData(null);
                      } else {
                        setFocusedFile(content);
                      }
                    }
                  }}
                  fileFocused={focusedFile && (focusedFile.Key === content.Key)}
                  setSelected={(siblings: any) => {
                    if (folderList[folder.Prefix].isSelected && fileList[content.Key].isSelected) {
                      setSingleFolder(folder.Prefix, false, null);
                    }
                    if (!fileList[content.Key].isSelected) {
                      let folderSelected = true;
                      siblings.files.forEach((file: File) => {
                        if(content.Key === file.Key) {
                          return;
                        }
                        if (!fileList[file.Key].isSelected) {
                          folderSelected = false;
                        }
                      })
                      siblings.folders.forEach((folder: Folder) => {
                        if (!folderList[folder.Prefix].isSelected) {
                          folderSelected = false;
                        }
                      })
                      if (folderSelected) {
                        setSingleFolder(folder.Prefix, true, null);
                      }
                    }
                    setSingleFile(content.Key, content, !fileList[content.Key].isSelected);
                  }}
                  fileSelected={fileList[content.Key].isSelected}
                  humanFileSize={humanFileSize}
                  s3Client={s3Client}
                  bucket={bucket}
                  refetchS3={() => refetchS3Local(folder.Prefix)}
                  siblings={{files, folders}}
                  removeFiles={removeFiles}
                  setIsLocked={setIsLocked}
                  namespace={namespace}
                  searchMode={searchMode}
                />
              )
              })
            }
          </div>
        )
      }
    </>
  );
}

export default Folder;
