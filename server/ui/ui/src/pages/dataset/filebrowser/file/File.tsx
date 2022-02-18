// vendor
import React, { FC, useState, useEffect, useRef } from 'react';
import { useDrag } from 'react-dnd'
import Moment from 'moment';
import classNames from 'classnames';
import { InputText } from 'Components/form/text/index';
import { CopyObjectCommand, DeleteObjectCommand } from "@aws-sdk/client-s3";
// components
import Actions from '../actions/Actions';
// css
import './File.scss';
/* eslint-disable */
const getClass = require('file-icons-js').getClass;
// typescript not supported with module, so require must be used instead
/* eslint-enable */

interface Props {
  file: {
    Key: string;
    Size: number;
    LastModified: any;
  };
  setFocus: any;
  fileFocused: boolean;
  setSelected: any;
  fileSelected: boolean;
  humanFileSize: any;
  s3Client: any;
  bucket: string;
  refetchS3: any;
  siblings: any;
  removeFiles: (fileKeys: Array<string>) => void;
  setIsLocked: (locked: boolean) => void;
  namespace: string;
  searchMode: boolean;
}


const File:FC<Props> = ({
  file,
  setFocus,
  fileFocused,
  setSelected,
  fileSelected,
  humanFileSize,
  s3Client,
  bucket,
  refetchS3,
  siblings,
  removeFiles,
  setIsLocked,
  namespace,
  searchMode,
}: Props) => {

  const [renameMode, setRenameMode] = useState(false);
  const inputRef = useRef(null);
  const checkboxRef = useRef(null);
  const nameRef = useRef(null);

  const splitName = file.Key.split('/');
  const fileName = splitName[splitName.length - 1];

  const [editedFileName, setFileName] = useState('');

  const path = splitName.slice(0, -1);
  const joinedPath = `${path.join('/')}/`;

  useEffect(() => {
    if (checkboxRef.current && !searchMode) {
      checkboxRef.current.style.marginLeft = `${(file.Key.split('/').length - 2) * 30}px`;
    } else if(nameRef.current) {
      nameRef.current.style.marginLeft = `${(file.Key.split('/').length - 2) * 30}px`;
    }
  }, [checkboxRef.current, nameRef.current])

  const [{ isDragging }, drag, preview] = useDrag(
    () => ({
      type: 'FILE',
      item: { file, fileName, sourcePath: joinedPath, siblings },
      collect: (monitor) => ({
        isDragging: !!monitor.isDragging(),
      }),
      end: (item, monitor) => {
        if (!monitor.didDrop()) {
          return
        } else {
          refetchS3();
        }
      },
    }),
    [],
  )

  const handleRename = () => {
    if (editedFileName === fileName) {
      setRenameMode(false);
      return;
    } else {
      const copyInput = {
          Bucket: bucket,
          CopySource: `${bucket}/${file.Key}`,
          Key: `${joinedPath}${editedFileName}`,
        }
        setIsLocked(true)
        const command = new CopyObjectCommand(copyInput);
        s3Client.send(command)
          .then((res: any) => {
            const deleteCommand = new DeleteObjectCommand({
              Key: file.Key,
              Bucket: bucket,
            })
            s3Client.send(deleteCommand)
            .then(() => {
              setIsLocked(false)
              refetchS3();
              removeFiles([file.Key]);
            });
          });
    }
  }

  useEffect(() => {
    if(renameMode) {
      inputRef.current.focus();
      inputRef.current.value = fileName;
    }
  }, [renameMode])

  const updateFileName = (evt: any) => {
    if((evt.key === 'Enter') && (editedFileName.length > 0)){
      handleRename();
    } else {
      setFileName(evt.target.value);
    }
  }

  const fileCSS = classNames({
    File: true,
    'File--focused': fileFocused,
  })

  const checkboxCSS = classNames({
    File__checkbox: true,
    'File__checkbox--selected': fileSelected,
  });

  return (
    <div className={fileCSS}
      onClick={setFocus}
      ref={searchMode ? null : drag}
    >
      {
        !searchMode && (
          <div
            ref={checkboxRef}
            className={checkboxCSS}
            onClick={(() => setSelected(siblings))}
          />
        )
      }
        {
          searchMode && <div className="File__checkbox--placeholder" />
        }
      <div className="File__file" ref={nameRef}>
        <div className={`File__icon ${getClass(fileName)}`} />
        {renameMode && (
          <InputText
            inputRef={inputRef}
            label=""
            placeholder={fileName}
            updateValue={updateFileName}
          />
        )}
        {!renameMode && fileName}
      </div>
      <div className="File__size">
        {humanFileSize(file.Size)}
      </div>
      <div className="File__modified">
        {Moment((file.LastModified)).fromNow()}
      </div>
      <div className="File__action">
        <Actions
          s3Client={s3Client}
          namespace={namespace}
          bucket={bucket}
          fileKey={file.Key}
          refetchS3={refetchS3}
          setRenameMode={() => {setRenameMode(!renameMode)}}
          handleRename={handleRename}
          renameMode={renameMode}
          removeFiles={removeFiles}
          allowRename={(editedFileName.length > 0)}
          searchMode={searchMode}
        />
      </div>
    </div>
  );
}

export default File;
