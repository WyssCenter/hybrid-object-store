// vendor
import React, { FC, useState } from 'react';
import { ListObjectsV2Command, GetObjectCommand, DeleteObjectCommand, GetBucketAclCommand } from "@aws-sdk/client-s3";
import ReactTooltip from 'react-tooltip';
import classNames from 'classnames';
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";
// environment
import { get } from 'Environment/createEnvironment';
// components
import { PrimaryButton, SecondaryButton } from 'Components/button/index';
import { ReactComponent as Download } from 'Images/icons/file-download.svg';
import { ReactComponent as Delete } from 'Images/icons/file-delete-active.svg';
import { ReactComponent as Rename } from 'Images/icons/file-edit-active.svg';
import { ReactComponent as Add } from 'Images/icons/file-new-folder-active.svg';

// css
import './Actions.scss';

interface File {
  Key: string;
  Size: number;
  LastModified: any;
}

interface Folder {
  Prefix: string;
}

interface Props {
  fileKey?: string;
  folderKey?: string;
  folder?: boolean;
  bucket: string;
  s3Client: any;
  refetchS3: any;
  renameMode: boolean;
  setRenameMode: () => void;
  removeFiles: (fileKeys: Array<string>) => void;
  removeFolders?: (folderKeys: Array<string>) => void;
  handleRename: () => void;
  allowRename: boolean;
  setNewFolderVisible?: (visible: boolean) => void;
  namespace?: string;
  searchMode: boolean;
  setErrorModal: any;
}


const Actions:FC<Props> = ({
  fileKey,
  folderKey,
  folder,
  bucket,
  s3Client,
  refetchS3,
  setRenameMode,
  renameMode,
  removeFiles,
  removeFolders,
  handleRename,
  allowRename,
  namespace,
  setNewFolderVisible,
  searchMode,
  setErrorModal,
}: Props) => {

  const [deleteMode, setDeleteMode] = useState(false);

  const deleteFile = () => {
    const command = new DeleteObjectCommand({
      Key: fileKey,
      Bucket: bucket,
    })
    try {
      s3Client.send(command)
      .then(() => {
        removeFiles([fileKey]);
        refetchS3()
      })
      .catch((error: any) => {
        console.log(error)
        setErrorModal({
          visible: true,
          action: 'deleting file',
          error: error.message,
        })
        setDeleteMode(false);
      })
    } catch (err) {
      console.log(err);
    }
  }

  const deleteFolder = () => {
    clearS3Folders(bucket, folderKey);
  }

  const clearS3Folders = async(bucket: string, source: string) => {
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
              setDeleteMode(false);
              setErrorModal({
                visible: true,
                action: 'deleting folder(s)',
                error: error.message,
              })
            })
          })
          } else {
            removeFolders([source]);
          }
          if (res.CommonPrefixes) {
            res.CommonPrefixes.forEach(async (folder: Folder) => {
              clearS3Folders(
                bucket,
                `${folder.Prefix}`,
              );
            })
          }
        })
    }

    const downloadFile = async() => {
      const command = new GetObjectCommand({
          Bucket: bucket,
          Key: fileKey
        })
      const url = await getSignedUrl(s3Client, command, { expiresIn: 3600 });
      window.open(url, '__blank');
    };
    const actionsCSS = classNames({
      Actions: true,
      'Actions--search': searchMode,
    });
    const downloadCSS = classNames({
      'Actions__Button Actions__Button--download': true,
      'Actions__Button--disabled': folder,
    })

  return (
    <div className={actionsCSS}>
      {
        folder && (
        <div
          data-tip="Not supported in search mode."
          role="presentation"
          data-tip-disable={!searchMode}
        >
          <button
            className="Actions__Button Actions__Button--add"
            onClick={() => {setNewFolderVisible(true)}}
            disabled={searchMode}
          >
            <Add />
          </button>
          {
            searchMode && (
            <ReactTooltip
              place="bottom"
              effect="solid"
            />
            )
          }
        </div>
        )
      }
      {
        renameMode && (
          <>
            <PrimaryButton
              click={handleRename}
              disabled={!allowRename}
              text="Rename"
            />
            <SecondaryButton
              click={() => {setRenameMode()}}
              text="Cancel"
            />
          </>
        )
      }
      {
        deleteMode && (
          <>
            <PrimaryButton
              click={folder ? deleteFolder : deleteFile}
              text="Delete"
            />
            <SecondaryButton
              click={() => {setDeleteMode(false)}}
              text="Cancel"
            />
          </>
        )
      }
      { !deleteMode && !renameMode && (
        <>
          <div
            data-tip="Directory downloads are not supported."
            role="presentation"
            data-tip-disable={!folder}
          >
            <button
              className={downloadCSS}
              onClick={() => {downloadFile()}}
              disabled={folder}
            >
              <Download />
            </button>
            {
              folder && searchMode && (
              <ReactTooltip
                place="bottom"
                effect="solid"
              />
              )
            }
          </div>
          <div
            data-tip="Not supported in search mode."
            role="presentation"
            data-tip-disable={!searchMode}
          >
            <button
              className="Actions__Button Actions__Button--rename"
              onClick={setRenameMode}
              disabled={searchMode}
            >
              <Rename />
            </button>
          {
            searchMode && (
            <ReactTooltip
              place="bottom"
              effect="solid"
            />
            )
          }
        </div>
          <div
            data-tip="Not supported in search mode."
            role="presentation"
            data-tip-disable={!searchMode}
          >
            <button
              className="Actions__Button Actions__Button--delete"
              onClick={() => { setDeleteMode(true) }}
              disabled={searchMode}
            >
              <Delete />
            </button>
          {
            searchMode && (
            <ReactTooltip
              place="bottom"
              effect="solid"
            />
            )
          }
        </div>
        </>
      )}
    </div>
  );
}

export default Actions;
