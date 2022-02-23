// vendor
import React, { FC, useState, useEffect } from 'react';
import { ListObjectsV2Command, DeleteObjectCommand } from "@aws-sdk/client-s3";
import { PrimaryButton, SecondaryButton } from 'Components/button/index';

// css
import './MultiAction.scss';


interface File {
  Key: string;
  Size: number;
  LastModified: any;
}

interface Folder {
  Prefix: string;
}

interface Props {
  folderList: any;
  fileList: any;
  s3Client: any;
  removeFiles: (fileKeys: Array<string>) => void;
  removeFolders?: (folderKeys: Array<string>) => void;
  bucket: string;
  setErrorModal: any;
}


const StaticButtons:FC<Props> = ({
  folderList,
  fileList,
  s3Client,
  removeFiles,
  removeFolders,
  bucket,
  setErrorModal,
}: Props) => {

  const [deleteMode, setDeleteMode] = useState(false);

  const filesSelected = Object.keys(fileList).map((fileKey) => fileList[fileKey]).filter((file: any) => file.isSelected).map((file: any) => file.file.Key);

  const foldersSelected = Object.keys(folderList).map((folderKey) => ({Prefix: folderKey, state: folderList[folderKey]})).filter((folder: any) => folder.state.isSelected).map((folder: any) => folder.Prefix);

  const selectedCountFiles = filesSelected.length;
  const selectedcountFolders = foldersSelected.length;

  const deleteSingleFile = (key: string) => {
    const command = new DeleteObjectCommand({
      Key: key,
      Bucket: bucket,
    })
    try {
      s3Client.send(command)
      .then(() => {
        removeFiles([key]);
      })
      .catch((error: any) => {
        if(setDeleteMode) {
          setDeleteMode(false);
        }
        setErrorModal({
          visible: true,
          action: 'deleting selected',
          error: error.message,
        })
      })
    } catch (err) {
      console.log(err);
    }
  }

  const deleteMultipleFiles = (keys: string[]) => {
    keys.forEach((key) => deleteSingleFile(key));
  }

  const deepClearFolder = async(bucket: string, source: string) => {
      if (!source.endsWith('/')) {
        return Promise.reject(new Error('source or dest must ends with fwd slash'));
      }
      let folderDeleted = false;
        s3Client.send(new ListObjectsV2Command({
          Bucket: bucket,
          Prefix: source,
          Delimiter: '/',
        }))
        .then((res: any) => {
          if (res.Contents) {
            res.Contents.forEach(async (file: File) => {
              s3Client.send(new DeleteObjectCommand({
                Key: file.Key,
                Bucket: bucket,
              }))
            .then(() => {
              removeFiles([file.Key]);
              if (!folderDeleted) {
                folderDeleted = true;
                removeFolders([source]);
              }
            })
            .catch((error: any) => {
              if(setDeleteMode) {
                setDeleteMode(false);
              }
              setErrorModal({
                visible: true,
                action: 'deleting selected',
                error: error.message,
              })
            })
          })
          }
          if (res.CommonPrefixes) {
            res.CommonPrefixes.forEach(async (folder: Folder) => {
              deepClearFolder(
                bucket,
                `${folder.Prefix}`,
              );
            })
          }
        })
    }

  const deleteMultipleFolders = (keys: string[]) => {
    keys.forEach((key) => deepClearFolder(bucket, key));
  }

  const handleDeleteAction = () => {
    deleteMultipleFiles(filesSelected);
    deleteMultipleFolders(foldersSelected);
  }



  const createString = (fileCount: number, folderCount: number) => {
    if(folderCount === 0 && (fileCount > 0)) {
      if (fileCount === 1) {
        return '1 file selected';
      }
      return `${fileCount} files selected`;
    }
    if (fileCount === 0 && (folderCount > 0)) {
      if (folderCount === 1) {
        return '1 folder selected';
      }
      return `${folderCount} folders selected`
    }
    if (fileCount === 1 && folderCount === 1) {
      return `1 folder and 1 file selected`
    }
    if (fileCount === 1) {
      return `${folderCount} folders and 1 file selected`
    }
    if (folderCount === 1) {
      return `1 folder and ${fileCount} files selected`
    }
    return `${folderCount} folders and ${fileCount} files selected `
  }
  const displayedString = createString(selectedCountFiles, selectedcountFolders)

  return (
    <div className="MultiAction flex">
      <div className="MultiAction--count">
        {displayedString}
      </div>
      <div className="MultiAction--actions">
        {
            deleteMode && (
              <>
                <SecondaryButton
                  click={handleDeleteAction}
                  text="Confirm"
                />
                <PrimaryButton
                  click={() => {setDeleteMode(false)}}
                  text="Cancel"
                />
              </>
            )
          }
          {
            !deleteMode && (
              <div
                className="MultiAction--delete"
                onClick={() => setDeleteMode(true)}
              >
                Delete
              </div>
            )
          }
      </div>
    </div>
  );
}

export default StaticButtons;
