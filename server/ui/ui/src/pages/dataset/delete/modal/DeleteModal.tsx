// vendor
import React, { FC, useState, useRef } from 'react';
// components
import Modal from 'Components/modal/Modal';
import { InputText } from 'Components/form/text/index';
import { WarningButton } from 'Components/button/index';
// css
import './DeleteModal.scss'



interface Props {
  datasetName: string;
  handleDeleteClick: () => void;
  hideModal: () => void;
  isVisible: boolean;
  deleteHours: number;
}



const DeleteModal: FC<Props> = (
  {
    datasetName,
    handleDeleteClick,
    hideModal,
    isVisible,
    deleteHours,
  }: Props
) => {
  const inputRef = useRef(null);
  const [confirmName, setConfirmName] = useState(null);
  const handleClick = () => {
    handleDeleteClick();
    inputRef.current.value = '';
    hideModal()
  }

  const updateConfirmName = (evt: any) => {
    setConfirmName(evt.target.value);
  }

  if (!isVisible) {
    return null;
  }
  return (
    <Modal
      handleClose={hideModal}
      header="Delete Dataset"
    >
      <section className="DeleteModal flex flex-direction--column justify--space-around">
        <div>
          <p>
            Are you sure? All data in this dataset will be lost.
          </p>
          <p>
            This action will permanently delete the dataset
            {' '}
            <b>{datasetName}</b>
            { ` in ${deleteHours} hours` }
            . If this is intended please type the name of the dataset to confirm and select Confirm Delete.
          </p>
        </div>
        <InputText
          inputRef={inputRef}
          label=""
          placeholder="Enter Dataset Name"
          updateValue={updateConfirmName}
        />
        <WarningButton
          click={handleClick}
          disabled={confirmName !== datasetName}
          text="Confirm Delete"
        />
      </section>
    </Modal>
  )
}


export default DeleteModal;
