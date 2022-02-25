// vendor
import React, { FC, useState, useEffect, useRef } from 'react';
import { Upload } from "@aws-sdk/lib-storage"
import { NativeTypes } from 'react-dnd-html5-backend'
import { useDrop } from 'react-dnd'
import { S3Client, ListObjectsV2Command, CopyObjectCommand, DeleteObjectCommand, HeadObjectCommand } from "@aws-sdk/client-s3";
import classNames from 'classnames';
// environment
import { get } from 'Environment/createEnvironment';
// provider
import { DndProvider } from 'react-dnd'
import { HTML5Backend } from 'react-dnd-html5-backend'
// components
import Modal from 'Components/modal/Modal';
import { PrimaryButton } from 'Components/button/index';
import File from './file/File';
import Details from './details/Details';
import {Folder} from './folder/Folder';
import TagsModal from './tagsModal/TagsModal';
import StaticButtons from './staticButtons/StaticButtons';
import MultiAction from './multiAction/MultiAction';
import ProgressBar from './progressBar/ProgressBar';
import NewFolder from './newFolder/NewFolder';
import Header from './header/Header';
// css
import { ReactComponent as Spinner } from 'Images/loaders/status-spinner.svg';
import './FileBrowser.scss';

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

interface Folder {
  Prefix: string;
}

interface File {
  Key: string;
  Size: number;
  LastModified: any;
}

interface Props {
  bucket: string;
  dataset: string;
  namespace: string;
  setDetailsVisible: (visible: boolean) => void;
  detailsVisible: boolean;
  lockBrowser: boolean;
}

const humanFileSize = (bytes: number) => {
  const si = true;
  const thresh = si ? 1000 : 1024;
  if (Math.abs(bytes) < thresh) {
    return `${bytes}B`;
  }
  const units = si
    ? ['KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
    : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
  let u = -1;
  do {
    bytes /= thresh;
    u += 1;
  } while (Math.abs(bytes) >= thresh && u < units.length - 1);
  return `${bytes.toFixed(1)} ${units[u]}`;
};

  const getS3 = async(
    accessKeyId: string,
    secretAccessKey: string,
    sessionToken: string,
    expiration: any,
    endpoint: string,
    region: string,
    Prefix: string,
    Bucket: string,
    isMinio: boolean,
    recountAttempt: number,
    callback: any,
    failCallback: any,
  ) => {
    try {
      const s3config = {
        credentials: {
          accessKeyId,
          secretAccessKey,
          sessionToken,
          expiration,
        },
        region: isMinio ? 'us-east-1':  region,
        endpoint,
        forcePathStyle: isMinio,
      }
      const ListObjectsInput = {
        Bucket,
        Delimiter: '/',
        Prefix,
      }

      const client = new S3Client(s3config);
      const command = new ListObjectsV2Command(ListObjectsInput);


      await client.send(command)
      .then(res => {
        callback(res, client)
      })
      .catch((err) => {
        if (recountAttempt > 5) {
          failCallback();
        } else {
          setTimeout(() => {
            getS3(
              accessKeyId,
              secretAccessKey,
              sessionToken,
              expiration,
              endpoint,
              region,
              Prefix,
              Bucket,
              isMinio,
              recountAttempt + 1,
              callback,
              failCallback,
            )
          }, 1000);
        }
      })
    } catch (error) {
      console.log('error', error)
    }
  }

const FileBrowser:FC<Props> = ({
  bucket,
  dataset,
  namespace,
  setDetailsVisible,
  detailsVisible,
  lockBrowser,
}: Props) => {
  const [canDropFolder, setCanDropFolder] = useState(null);
  const [isMinio, setIsMinio] = useState(false);
  const [newFolderVisible, setNewFolderVisible] = useState(null);
  const [focusedFile, setFocusedFile] = useState(null);
  const [focusedFileData, setFocusedFileData] = useState(null);
  const [s3Client, setS3Client] = useState(null);
  const [preparedUpload, setPreparedupload] = useState(null)
  const [tagModalVisible, setTagModalVisible] = useState(false);
  const [pendingUpload, setPendingUpload] = useState(null);
  const [selectedSort, setSelectedSort] = useState({
    type: 'file',
    ascending: false,
  });
  const initialFiles: FileObject = {};
  const [fileList, setFileList] = useState(initialFiles);
  const initialFolders: FolderObject  = {};
  const [folderList, setFolderList] = useState(initialFolders);
  const [isLocked, setIsLocked] = useState(lockBrowser);
  const [searchData, setSearchData] = useState(null);
  const [loadingSearch, setLoadingSearch] = useState(false);
  const [errorFetch, setErrorFetch] = useState(false);
  const [loadingFetch, setLoadingFetch] = useState(true);
  const [errorModal, setErrorModal] = useState({
    visible: false,
    action: '',
    error: '',
  })
  const [uploadData, setUploadData] = useState({
    uploading: false,
    totalFiles: null,
    totalSize: null,
    finishedFiles: null,
    finishedSize: null,
    pct: null,
  });

  const effectiveFileList = searchData
    ? searchData.fileList
    : fileList;

  const effectiveFolderList = searchData
    ? searchData.folderList
    : folderList;

  const setSearchFolderExpanded = (folderName: string, isExpanded: boolean) => {
    setSearchData((prevData: any) => {
      const newFolderList = Object.assign({}, prevData.folderList);
      newFolderList[folderName].isExpanded = isExpanded;
      return {
        fileList: prevData.fileList,
        folderList: newFolderList,
        size: prevData.size,
      }
    })
  }


  const formatSearchData = (files: any[]) => {
    const folderObject: any = {}
    const fileObject: any = {}
    files.forEach((file: any) => {
      const splitPath = file.file_path.split('/');
      if (splitPath.length > 1) {
        folderObject[`${dataset}/${splitPath.slice(0, -1).join('/')}/`] = {isExpaned: false};
      }
      const fileKey = `${dataset}/${file.file_path}`;
      fileObject[fileKey] = {
        file: {
          Key: fileKey,
          LastModified: file.last_modified_date,
          Size: file.size_bytes,
        }
      }
    })
    setSearchData({
      fileList: fileObject,
      folderList: folderObject,
      size: files.length
    })
  }

  const searchFiles = (metadata: any, modifiedBefore: string, modifiedAfter: string) => {
    const formattedData = metadata ? Object.keys(metadata).map((key) => {
      return `${key}:${metadata[key]}`
    }) : [];
    const finalFormat = `&metadata=${formattedData.join(',')}`;
    const metadataParsed = (metadata && (Object.keys(metadata).length > 0)) ? finalFormat : '';
    const beforeModifiedParsed = modifiedBefore ? `&modified_before=${new Date(modifiedBefore).toISOString()}` : '';
    const afterModifiedParsed = modifiedAfter ? `&modified_after=${new Date(modifiedAfter).toISOString()}` : '';
    setIsLocked(true);
    setLoadingSearch(true);
    get(`search?namespace=${namespace}&dataset=${dataset}${metadataParsed}${beforeModifiedParsed}${afterModifiedParsed}`)
    .then(res => res.json())
    .then(data => {
      setIsLocked(false);
      setLoadingSearch(false);
      formatSearchData(data.results);
    })
  }


  useEffect(() => {
    if(focusedFile && !detailsVisible){
      setDetailsVisible(true)
    } else if (!focusedFile && detailsVisible) {
      setDetailsVisible(false)
    }
  }, [focusedFile])

  const s3Ref = useRef(null);
  const uploadRef = useRef(null);

  s3Ref.current = s3Client;
  uploadRef.current = uploadData.uploading;

  const batchUpload = (files: any, path: string) => {
    if (!uploadRef.current) {
      setTagModalVisible(true);
      setPendingUpload({
        files,
        path,
      })
    }
  }

  const upload = (files: any, path: string, metaObject: any) => {
    files = files.filter((file: any) => {
      if (file.file && (file.file.name === '.DS_Store')) {
        return false;
      } else if ((file.name === '.DS_Store')) {
        return false;
      }
      return true;
    })
    setIsLocked(true);
    const totalFiles = files.length;
    let totalSize = 0;
    let finishedFiles = 0;
    const downloadProgress: any = {};
    files.forEach((file: any) => totalSize += file.file ? file.file.size : file.size);
    setUploadData((prevData) =>({
      ...prevData,
      uploading: true,
      totalFiles,
      totalSize,
      finishedFiles,
      finishedSize: 0,
      pct: 0,
    }))


    files.forEach((file: any) => {
      const fileName = file.file ? file.file.name : file.name;
      const key = file.file
        ? `${path}${file.fullPath.slice(1)}/${fileName}`
        :`${path}${fileName}`;
      const target = {
        Bucket: bucket,
        Key: key,
        Body: file.file ? file.file : file,
        Metadata: metaObject || {},
      }

      try {
        const multiUpload = new Upload({
          client: s3Ref.current,
          leavePartsOnError: false,
          params: target,
        })

        multiUpload.on('httpUploadProgress', (progress) => {
          downloadProgress[key] = progress.loaded;
          if (progress.loaded === progress.total) {
            finishedFiles += 1;
          }
          const totalProgress = Object.values(downloadProgress).reduce((a: number,b: number)=>a+b);
          const downloaded = totalProgress as number;
          const pct: any = ((downloaded / totalSize) * 100).toFixed();
          setUploadData((prevData) =>({
            ...prevData,
            uploading: true,
            totalFiles,
            totalSize,
            finishedFiles,
            finishedSize: downloaded,
            pct,
          }))
        })
        multiUpload.done()
        .then(() => {
          refetchS3(path)
          if (totalFiles === finishedFiles) {
            setUploadData((prevData) =>({
              ...prevData,
              uploading: false,
              totalFiles: null,
              totalSize: null,
              finishedFiles: null,
              downloaded: null,
              pct: null,
            }))
            setIsLocked(false)
          }
        })
        .catch(error => {
          if (!errorModal.visible) {
            setErrorModal({
              visible: true,
              action: 'uploading file(s)',
              error: error.message,
            })
          }
          setUploadData((prevData) =>({
            ...prevData,
            uploading: false,
            totalFiles: null,
            totalSize: null,
            finishedFiles: null,
            downloaded: null,
            pct: null,
          }))
          setIsLocked(false)
        })
      } catch (err) {
        console.log(err)
      }
    })
  };

  const s3CopyFolderHandler = (bucket: string, source: string, dest: string) => {
    setIsLocked(true);
    let pendingFolders = 0;
    let pendingFiles = 0;

    const s3CopyFolder = async(bucket: string, source: string, dest: string) => {
      if (!source.endsWith('/') || !dest.endsWith('/')) {
        return Promise.reject(new Error('source or dest must ends with fwd slash'));
      }
        s3Ref.current.send(new ListObjectsV2Command({
          Bucket: bucket,
          Prefix: source,
          Delimiter: '/',
        }))
        .then((res: any) => {
          pendingFolders -= 1;
          if (res.CommonPrefixes) {
            pendingFolders += res.CommonPrefixes.length;
            res.CommonPrefixes.forEach(async (folder: Folder) => {
              s3CopyFolder(
                bucket,
                `${folder.Prefix}`,
                `${dest}${folder.Prefix.replace(res.Prefix, '')}`,
              );
            })
          }
          if (res.Contents) {
            const contentCount = res.Contents.length;
            pendingFiles += contentCount;
            let resolvedCount = 0;
            res.Contents.forEach(async (file: File) => {
            s3Ref.current.send(new CopyObjectCommand({
              Bucket: bucket,
              CopySource: `${bucket}/${file.Key}`,
              Key: `${dest}${file.Key.replace(res.Prefix, '')}`,
            }))
            .then(() => {
              s3Ref.current.send(new DeleteObjectCommand({
                Key: file.Key,
                Bucket: bucket,
              }))
              .then(() => {
                removeFiles([file.Key]);
                pendingFiles -= 1;
                resolvedCount +=1 ;
                if (resolvedCount === contentCount) {
                  s3Ref.current.send(new ListObjectsV2Command({
                    Bucket: bucket,
                    Delimiter: '/',
                    Prefix: dest,
                  }))
                  .then((res: any) => {
                    const newFilesObject: FileArray = {};
                    const fileKeys = (res.Contents && res.Contents.map((content: any) => {
                      newFilesObject[content.Key] = content;
                      return content.Key;
                    })) || [];
                    setMultipleFiles(fileKeys, newFilesObject, null)
                    const folderKeys = (res.CommonPrefixes && res.CommonPrefixes.map((content: any) => content.Prefix)) || [];
                    setMultipleFolders(folderKeys, null, null);
                    setSingleFolder(dest, false, false)
                    removeFolders([source]);
                    if (pendingFolders === 0 && pendingFiles === 0){
                      setIsLocked(false);
                    }
                  })
                }
              })
            })
            .catch((error: any) => {
              setIsLocked(false);
              setErrorModal({
                visible: true,
                action: 'modifying folder(s)',
                error: error.message,
              })
            })
          })
          } else {
            removeFolders([source]);
          }
          if (pendingFolders === 0 && pendingFiles === 0){
            setIsLocked(false);
          }
        })
    }

    pendingFolders = 1;
    s3CopyFolder(bucket, source, dest);
  }

  const folderReader = (folder: any, fileArray: any[], path: string) => {
    setPreparedupload((prevState: any) => ({
      ...prevState,
      fileArray,
      openCount: prevState ? prevState.openCount + 1 : 1,
      closedCount: prevState ? prevState.closedCount : 0,
      path,
    }));

    const directoryReader = folder.createReader();
    directoryReader.readEntries((entries: any) => {
      const totalEntries = entries.length;
      let parsedEntries = 0;
      if (totalEntries === 0) {
        setPreparedupload((prevState: any) => ({
          ...prevState,
          closedCount: prevState.closedCount + 1,
        }));
      }
      entries.forEach((entry: any) => {
        if (entry.isDirectory) {
          folderReader(entry, fileArray, path);
          parsedEntries +=1;
          if (parsedEntries === totalEntries) {
            setPreparedupload((prevState: any) => ({
              ...prevState,
              closedCount: prevState.closedCount + 1,
            }));
          }
        } else {
          entry.file((file: any) => {
            fileArray.push({file, fullPath: folder.fullPath});
            parsedEntries +=1;
            if (parsedEntries === totalEntries) {
              setPreparedupload((prevState: any) => ({
                ...prevState,
                fileArray,
                openCount: prevState.openCount,
                closedCount: prevState.closedCount + 1,
                path,
              }));
            }
          })
        }
      })
    })
  }

  useEffect(() => {
    if (!preparedUpload) {
      return;
    }
    if (preparedUpload.openCount === preparedUpload.closedCount) {
      batchUpload(preparedUpload.fileArray, preparedUpload.path);
      setPreparedupload(null);
    }
  }, [preparedUpload])

  const [{ isOver, canDrop }, drop] = useDrop(
    () => ({
      accept: [NativeTypes.FILE, 'FILE'],
      drop: (item: any, monitor) => {
        const didDrop = monitor.didDrop()
        if (didDrop) {
          return
        }
        if
        ((item.sourcePath === `${dataset}/`) || (item.parentFolder === `${dataset}/`)
        ) {
          return;
        }
        if (uploadData.uploading) {return;}
        if (item.dataTransfer && item.dataTransfer.items) {
          const { items } = item.dataTransfer;
          const fileArray: any[] = [];
          Array.from(items).forEach((item: any) =>{
            setIsLocked(true);
            const folder = item.webkitGetAsEntry();
            if (folder && folder.isDirectory) {
              folderReader(folder, fileArray, `${dataset}/`)
            } else {
              fileArray.push(item.getAsFile());
            }
          })
          if (fileArray.length === items.length) {
            batchUpload(fileArray, `${dataset}/`)
          }
          return;
        }
        if (item.file) {
          setIsLocked(true);
          const copyInput = {
            Bucket: bucket,
            CopySource: `${bucket}/${item.file.Key}`,
            Key: `${dataset}/${item.fileName}`,
          }
          const command = new CopyObjectCommand(copyInput);
          s3Ref.current.send(command)
            .then((res: any) => {
              const deleteCommand = new DeleteObjectCommand({
                Key: item.file.Key,
                Bucket: bucket,
              })
              s3Ref.current.send(deleteCommand)
              .then(() => {
                refetchS3(`${dataset}/`);
                removeFiles([item.file.Key]);
                setIsLocked(false);
              });
            })
            .catch((error: any) => {
              setIsLocked(false);
              setErrorModal({
                visible: true,
                action: 'moving file(s)',
                error: error.message,
              })
            })
        }
        if (item.folder
          && item.parentFolder !== `${dataset}/`
          && item.folder.Prefix !== `${dataset}/`) {
          s3CopyFolderHandler(bucket, item.folder.Prefix, `${dataset}/${item.folderName}/`);
        }
        if (item.files) {
          batchUpload(item.files, `${dataset}/`);
        }
      },
      collect: (monitor: any) => {
        return {
          isOver: !!monitor.isOver(),
          canDrop: !!monitor.canDrop()
            &&  (monitor.getItem().sourcePath !== `${dataset}/`)
            &&  (monitor.getItem().parentFolder !== `${dataset}/`)
        }
      },
    }),
    [],
  )

  const setSingleFile = (fileKey: string, file: File, isSelected?: boolean) => {
    setFileList((prevState) => ({
      ...prevState,
      [fileKey]: {
        file,
        isSelected: (isSelected === null) ? prevState[fileKey].isSelected : isSelected
        },
    }))
  };

  const setMultipleFiles = (fileKeys: Array<string>, files: FileArray, isSelected: boolean) => {
    setFileList((prevState) =>
      {
        const newFiles: FileObject = {};
        fileKeys.forEach(fileKey => {
          newFiles[fileKey] = {
            file: files[fileKey],
            isSelected: prevState[fileKey] ?
              (isSelected === null) ? prevState[fileKey].isSelected : isSelected
              : (isSelected === null) ? false : isSelected,
          };
        })
        return {
          ...prevState,
          ...newFiles,
        }
    })
  }

  const removeFiles = (fileKeys: Array<string>) => {
    setFileList((prevState) => {
      const newFiles: FileObject = Object.assign({}, prevState);
      fileKeys.forEach(fileKey => {
        delete newFiles[fileKey];
      })
      return ({
        ...newFiles
      });
    })
  }

  const setSingleFolder = (folderKey: string, isSelected: boolean, isExpanded: boolean) => {
    setFolderList((prevState) => ({
      ...prevState,
      [folderKey]:  {
        isSelected: folderList[folderKey] ?
          (isSelected === null) ? folderList[folderKey].isSelected : isSelected
          : (isSelected === null) ? false : isSelected,
        isExpanded: folderList[folderKey] ?
          (isExpanded === null) ? folderList[folderKey].isExpanded : isExpanded
          : (isExpanded === null) ? false : isExpanded,
        },
    }))
  };

  const setMultipleFolders = (folderKeys: Array<string>, isSelected: boolean, isExpanded: boolean) => {

    setFolderList((prevState) => {
      const newFolders: FolderObject = {};

      folderKeys.forEach(folderKey => {
        newFolders[folderKey] = {
          isSelected: prevState[folderKey] ?
            (isSelected === null) ? prevState[folderKey].isSelected : isSelected
            : (isSelected === null) ? false : isSelected,
          isExpanded: prevState[folderKey] ?
            (isExpanded === null) ? prevState[folderKey].isExpanded : isExpanded
            : (isExpanded === null) ? false : isExpanded,
          };
      })
      return {
        ...prevState,
        ...newFolders,
      }
    })
  }

  const removeFolders = (folderKeys: Array<string>) => {
    setFolderList((prevState) => {
      const newFolders: FolderObject = Object.assign({}, prevState);
      folderKeys.forEach(folderKey => {
        delete newFolders[folderKey];
      })
      return ({
        ...newFolders
      });
    })
  }

  const setAllSelected = (isSelected: boolean) => {
    const folders = Object.keys(folderList);
    setMultipleFolders(folders, isSelected, null);

    const fileKeys = Object.keys(fileList);
    const files: FileArray  = {};
    fileKeys.forEach((fileKey) => files[fileKey] = fileList[fileKey].file);
    setMultipleFiles(fileKeys, files, isSelected);

  }

  let allSelected = Object.values(fileList).length > 0 || Object.values(folderList).length > 0;
  let someSelected = false;

  Object.keys(fileList).forEach((fileKey) => {
    if (!fileList[fileKey].isSelected) {
      allSelected = false;
    } else {
      someSelected = true;
    }
  });
  Object.keys(folderList).forEach((folderKey) => {
    if (!folderList[folderKey].isSelected) {
      allSelected = false;
    } else {
      someSelected = true;
    }
  });

  const refetchS3 = (Prefix: string) => {
    const ListObjectsInput = {
        Bucket: bucket,
        Delimiter: '/',
        Prefix,
      }
    const command = new ListObjectsV2Command(ListObjectsInput);
    s3Ref.current.send(command)
    .then((res: any) => {
      const newFilesObject: FileArray = {};
      const fileKeys = (res.Contents && res.Contents.map((content: any) => {
        newFilesObject[content.Key] = content;
        return content.Key;
      })) || [];
      setMultipleFiles(fileKeys, newFilesObject, null)
      const folderKeys = (res.CommonPrefixes && res.CommonPrefixes.map((content: any) => content.Prefix)) || [];
      setMultipleFolders(folderKeys, null, null);
    })
  }

  useEffect(() => {
    if (focusedFile) {
      const command = new HeadObjectCommand({
        Bucket: bucket,
        Key: focusedFile.Key
      })
      s3Client.send(command)
      .then((res: any) => {
        setFocusedFileData((prevData: any) => {
          if (prevData && !prevData.ETag) {
            return {
            ...res,
            Metadata: prevData.Metadata,
            }
          }
          return {
            ...res,
          }
        });
      })
      if (!isMinio) {
        get(`search/namespace/${namespace}/dataset/${dataset}/metadata?objectKey=${encodeURIComponent(focusedFile.Key)}`)
        .then(res => res.json())
        .then(data => {
          if (data && data.metadata) {
            setFocusedFileData((prevData: any) => {
              if (prevData && (prevData.ETag === focusedFile.ETag)) {
                return {
                  ...prevData,
                  Metadata: data.metadata
                }
              }
              return {
                Metadata: data.metadata,
              }
            })
          }
        })
      }
    } else if (focusedFile && focusedFileData) {
      setFocusedFileData(null);
    }
  }, [focusedFile])

  useEffect(() => {
    if (bucket) {
      setIsLocked(true);
      Promise.all([
        get(`namespace/${namespace}/sts`),
        get(`namespace/${namespace}`),
      ])
      .then(responses => Promise.all(responses.map(response => response.json())))
      .then(([data, namespaceData]) => {
        setIsMinio(namespaceData.object_store.type === 'minio');
        getS3(
          data.access_key_id,
          data.secret_access_key,
          data.session_token,
          data.expiration,
          data.endpoint,
          data.region,
          `${dataset}/`,
          bucket,
          namespaceData.object_store.type === 'minio',
          0,
          (res: any, client: any) => {
            setS3Client(client);
            const newFilesObject: FileArray = {};
            const fileKeys = (res.Contents && res.Contents.map((content: any) => {
              newFilesObject[content.Key] = content;
              return content.Key;
            })) || [];
            setMultipleFiles(fileKeys, newFilesObject, false);
            const folderKeys = (res.CommonPrefixes && res.CommonPrefixes.map((content: any) => content.Prefix)) || [];
            setMultipleFolders(folderKeys, false, false);
            setIsLocked(false);
            setLoadingFetch(false);
          },
          () => {
            setErrorFetch(true);
            setIsLocked(false);
          }
        );
      })
      .catch((error: Error) => {
          setIsLocked(false);
          setErrorFetch(true);
          console.log(error);
      });
    }
  }, [bucket, dataset, namespace])

  let files = Object.keys(effectiveFileList).map((fileKey) => effectiveFileList[fileKey].file).filter((file) => (file.Key.split('/').length === 2) && (file.Key !== `${dataset}/.dataset.yaml`));
  let folders = Object.keys(effectiveFolderList).map((folderKey) => ({ Prefix: folderKey })).filter((folder) => folder.Prefix.split('/').length === 3);

  const fileBrowserCSS = classNames({
    FileBrowser: true,
    'Folder--hovered': isOver
      && canDrop
      && !isLocked
      && (!canDropFolder || (canDropFolder && (Object.values(canDropFolder).indexOf(true) === -1)))
  })

  const sort = (files: any, isFolder: boolean) => {
    files.sort((a: any, b: any) => {
      // sort on name
      if (selectedSort.type === 'file') {
        const index = isFolder ? 2 : 1;
        const key = isFolder ? 'Prefix' : 'Key';
        const splitA = a[key].split('/');
        const nameA = splitA[splitA.length - index].toLowerCase();
        const splitB = b[key].split('/');
        const nameB = splitB[splitB.length - index].toLowerCase();
        if (selectedSort.ascending) {
          if (nameA < nameB) { return 1; }
          if (nameA > nameB) { return -1; }
          return 0;
        }
        if (nameA < nameB) { return -1; }
        if (nameA > nameB) { return 1; }
        return 0;
      }
      if (!isFolder) {
        // sort on size
        if (selectedSort.type === 'size') {
          if (selectedSort.ascending) {
            if (a.Size < b.Size) { return 1; }
            if (a.Size > b.Size) { return -1; }
            return 0;
          }
          if (a.Size < b.Size) { return -1; }
          if (a.Size > b.Size) { return 1; }
          return 0;
        }
        // sort on modified
        if (selectedSort.type === 'modified') {
          if (selectedSort.ascending) {
            if (a.LastModified < b.LastModified) { return 1; }
            if (a.LastModified > b.LastModified) { return -1; }
            return 0;
          }
          if (a.LastModified < b.LastModified) { return -1; }
          if (a.LastModified > b.LastModified) { return 1; }
          return 0;
        }
      }
    })
    return files;
  }

  files = sort(files, false);
  folders = sort(folders, true);

  if(errorFetch) {
    return (
      <div className="FileBrowser__Error">
        <p>
          There was an issue fetching files for this dataset. Refresh this page or try again later.
        </p>
        <p>
          If this issue persists please contact an administrator.
        </p>
      </div>
    )
  }

  return (
    <div
      className={fileBrowserCSS}
      ref={!searchData ? drop : null}
    >
      {
        isLocked && (
        <div
          className="FileBrowser--mask"
        />
        )
      }
      {
        errorModal.visible && (
          <Modal
            size="flex"
            header={`Error ${errorModal.action}`}
            handleClose={() => setErrorModal({
              visible: false,
              action: '',
              error: '',
            })}
          >
            <p className="">{
              errorModal.error === 'Access Denied.'
               ? 'Access Denied. This is likely due to a lack of permissions to modify the dataset. Read & Write permissions are required. '
               : errorModal.error
            }</p>
            <PrimaryButton
                click={() => setErrorModal({
                visible: false,
                action: '',
                error: '',
              })}
              text="Close"
            />
          </Modal>
        )
      }
      <TagsModal
        hideModal={() => setTagModalVisible(false)}
        isVisible={tagModalVisible}
        setPendingUpload={setPendingUpload}
        upload={upload}
        pendingUpload={pendingUpload}
        setIsLocked={setIsLocked}
        namespace={namespace}
        dataset={dataset}
      />
      <Details
        setDetailsVisible={setDetailsVisible}
        detailsVisible={detailsVisible}
        isVisible={focusedFileData}
        data={focusedFileData}
        focusedFile={focusedFile}
        setFocus={setFocusedFile}
        humanFileSize={humanFileSize}
        setFocusedFileData={setFocusedFileData}
        dataset={dataset}
        namespace={namespace}
      />
      {
        someSelected && !uploadData.uploading && (
          <MultiAction
            folderList={folderList}
            fileList={fileList}
            s3Client={s3Client}
            removeFiles={removeFiles}
            removeFolders={removeFolders}
            bucket={bucket}
            setErrorModal={setErrorModal}
          />
        )
      }
      {
        uploadData.uploading && !someSelected && (
          <ProgressBar
            uploadData={uploadData}
          />
        )
      }
      {
        !uploadData.uploading && !someSelected && (
          <StaticButtons
            s3Client={s3Client}
            dataset={dataset}
            namespace={namespace}
            bucket={bucket}
            refetchS3={refetchS3}
            batchUpload={batchUpload}
            setNewFolderVisible={setNewFolderVisible}
            searchFiles={searchFiles}
            searchData={searchData}
            clearSearchData={() => {setSearchData(null)}}
            loadingSearch={loadingSearch}
          />
        )
      }
      <Header
        allSelected={allSelected}
        setAllSelected={() => setAllSelected(!allSelected)}
        someSelected={someSelected}
        setSelectedSort={setSelectedSort}
        selectedSort={selectedSort}
        searchMode={!!searchData}
      />
      {
        newFolderVisible && (
          <NewFolder
            closeFolder={() => setNewFolderVisible(false)}
            isVisible={newFolderVisible}
            setSingleFolder={setSingleFolder}
            prefix={`${dataset}/`}
          />
        )
      }
      {
        folders && folders.map((content: Folder) => {
         return(
          <Folder
            key={content.Prefix}
            folder={content}
            setSelected={(folderExpanded: boolean, folderContents: FolderContents) => {
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
            folderSelected={folderList[content.Prefix] && folderList[content.Prefix].isSelected}
            setExpanded={(evt: Event) => {
              if(searchData) {
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
            folderExpanded={effectiveFolderList[content.Prefix].isExpanded}
            s3Client={s3Client}
            bucket={bucket}
            setSingleFolder={setSingleFolder}
            setMultipleFolders={setMultipleFolders}
            setSingleFile={setSingleFile}
            setMultipleFiles={setMultipleFiles}
            folderList={effectiveFolderList}
            fileList={effectiveFileList}
            setFocusedFile={setFocusedFile}
            focusedFile={focusedFile}
            humanFileSize={humanFileSize}
            removeFiles={removeFiles}
            removeFolders={removeFolders}
            refetchS3={refetchS3}
            siblings={null}
            depth={3}
            setCanDropFolder={setCanDropFolder}
            canDropFolder={canDropFolder}
            s3CopyFolder={s3CopyFolderHandler}
            setFocusedFileData={setFocusedFileData}
            batchUpload={batchUpload}
            sort={sort}
            isLocked={isLocked}
            setIsLocked={setIsLocked}
            folderReader={folderReader}
            namespace={namespace}
            searchMode={!!searchData}
            setSearchFolderExpanded={setSearchFolderExpanded}
            effectiveFolderList={effectiveFolderList}
            setErrorModal={setErrorModal}
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
              setSingleFile(content.Key, content, !fileList[content.Key].isSelected);
            }}
            fileSelected={fileList[content.Key] && fileList[content.Key].isSelected}
            humanFileSize={humanFileSize}
            s3Client={s3Client}
            bucket={bucket}
            refetchS3={() => refetchS3(`${dataset}/`)}
            siblings={{files, folders}}
            removeFiles={removeFiles}
            setIsLocked={setIsLocked}
            namespace={namespace}
            searchMode={!!searchData}
            setErrorModal={setErrorModal}
          />
         )
        })
      }
      {
        loadingFetch && (
          <div className="FileBrowser__loading">
            <Spinner />
            <p>
              Fetching Dataset files, please wait...
            </p>
          </div>
        )
      }
      {
        (!newFolderVisible
        && searchData === null
        && !loadingFetch
        && !errorFetch
        && folders.length === 0
        && files.length === 0) && (
          <div className="FileBrowser__empty">
            <div className="FileBrowser__empty--text">
              Drag and drop files here
            </div>
            <div className="FileBrowser__empty--divider">
              ———— or ————
            </div>
            <label
              htmlFor="add_file"
              className="flex justify--center"
            >
              <div
                className="Button"
              >
                Choose File...
              </div>
              <input
              id="add_file"
              className="hidden"
              type="file"
              onChange={(evt) => {
                const { files } = evt.target;
                batchUpload([files[0]], `${dataset}/`)
              }}
              />
            </label>
          </div>
        )
      }
    </div>
  );
}

const FileBrowserWrapper = (props: any) => (
    <DndProvider backend={HTML5Backend}>
      <FileBrowser {...props}/>
    </DndProvider>
)

export default FileBrowserWrapper;
