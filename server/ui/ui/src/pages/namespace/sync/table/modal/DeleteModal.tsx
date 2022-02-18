// vendor
import React, { FC, useState, useRef } from 'react';
// components
import Modal from 'Components/modal/Modal';
import { InputText } from 'Components/form/text/index';
import { WarningButton } from 'Components/button/index';
// css
import './DeleteModal.scss'



interface Props {
  namespaceName: string;
  handleDeleteClick: () => void;
  hideModal: () => void;
  isVisible: boolean;
}



const DeleteModal: FC<Props> = (
  {
    namespaceName,
    handleDeleteClick,
    hideModal,
    isVisible,
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
      header="Disable Sync"
    >
      <section className="DeleteModal flex flex-direction--column justify--space-around">
        <div>
          <p>
            Are you sure? You should ensure that no datasets in this namespace have syncing enabled before disabling namespace sync support. This will stop all datasets from syncing.
          </p>
          <p>
            This action will disable sync for
            {' '}
            <b>{namespaceName}</b>
            . If this is intended please type the name of the namespace to confirm and select Confirm Disable.
          </p>
        </div>
        <InputText
          inputRef={inputRef}
          label=""
          placeholder="Enter Namespace Name"
          updateValue={updateConfirmName}
        />
        <WarningButton
          click={handleClick}
          disabled={confirmName !== namespaceName}
          text="Confirm Disable"
        />
      </section>
    </Modal>
  )
}


export default DeleteModal;
