// vendor
import React, { FC, useState, useRef, useEffect } from 'react';
// components
import Modal from 'Components/modal/Modal';
import { InputText } from 'Components/form/text/index';
import { ReactComponent as FolderIcon } from 'Images/icons/folder.svg';
import { PrimaryButton, SecondaryButton } from 'Components/button/index';
// css
import './NewFolder.scss'



interface Props {
  isVisible: boolean;
  closeFolder: () => void;
  setSingleFolder: (folderKey: string, isSelected: boolean, isExpanded: boolean) => void;
  prefix: string;
}



const DeleteModal: FC<Props> = (
  {
    isVisible,
    closeFolder,
    setSingleFolder,
    prefix,
  }: Props
) => {
  const inputRef = useRef(null);
  const [folderName, setFolderName] = useState('');
  const handleClick = () => {
    setSingleFolder(`${prefix}${folderName}/`, false, false)
    inputRef.current.value = '';
    closeFolder();
  }

  const updateFolderName = (evt: any) => {
    if((evt.key === 'Enter') && (folderName.length > 0)){
      handleClick();
    } else {
      setFolderName(evt.target.value);
    }
  }

  useEffect(() => {
    if(isVisible) {
      inputRef.current.focus();
    }
  }, [isVisible])

  if (!isVisible) {
    return null;
  }
  return (
    <div
      className="Folder"
    >
        <div
          className="Folder__checkbox--placeholder"
        />
        <div className="Folder__file">
          <FolderIcon />
          <InputText
            inputRef={inputRef}
            label=""
            placeholder={folderName}
            updateValue={updateFolderName}
          />
        </div>
        <div className="Folder__size">
          {' '}
        </div>
        <div className="Folder__modified">
          {' '}
        </div>
        <div className="Folder__action Actions">
          <>
            <PrimaryButton
              click={handleClick}
              disabled={(folderName.length === 0)}
              text="Create"
            />
            <SecondaryButton
              click={closeFolder}
              text="Cancel"
            />
          </>
        </div>
    </div>
  )
}


export default DeleteModal;
