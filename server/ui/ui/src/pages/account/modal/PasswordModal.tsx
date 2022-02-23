// vendor
import React, { FC, useState, useRef } from 'react';
// environment
import { put } from 'Environment/createEnvironment';
// components
import Modal from 'Components/modal/Modal';
import { InputText } from 'Components/form/text/index';
import { PrimaryButton } from 'Components/button/index';
// css
import './PasswordModal.scss'



interface Props {
  hideModal: () => void;
  isVisible: boolean;
}

const PasswordModal: FC<Props> = (
  {
    hideModal,
    isVisible,
  }: Props
) => {
  const inputRef = useRef(null);
  const [confirmName, setConfirmName] = useState('');
  const [newName, setNewName] = useState('');
  const [newNameConfirm, setNewNameConfirm] = useState('');
  const [showRegexWarning, setRegexWarning] = useState(false);
  const [showPasswordMismatch, setShowPasswordMismatch] = useState(false);
  const [error, setError] = useState(null);
  const [succesfullyChanged, setSuccessfullyChanged] = useState(false);
  const [fetching, setFetching] = useState(false);

  const updateConfirmName = (evt: any) => {
    setConfirmName(evt.target.value);
    if((evt.key === 'Enter') && (!disableChange)){
      handlePasswordConfirm();
    }
  }
  const updateNewName = (evt: any) => {
    setNewName(evt.target.value);
    if((evt.key === 'Enter') && (!disableChange)){
      handlePasswordConfirm();
    }
  }
  const updateNewNameConfirm = (evt: any) => {
    setNewNameConfirm(evt.target.value);
    if((evt.key === 'Enter') && (!disableChange)){
      handlePasswordConfirm();
    }
  }

  const handlePasswordConfirm = () => {
    if (newNameConfirm !== newName) {
      setShowPasswordMismatch(true);
      setTimeout(() => {
        setShowPasswordMismatch(false);
      }, 5000)
      return;
    }
    setFetching(true);
    put(
      'password',
      {
        current: confirmName,
        new: newName,
      },
      true,
    )
    .then(response => {
      setFetching(false);
      if (response.status === 204) {
        setSuccessfullyChanged(true);
        setTimeout(() => {
          setSuccessfullyChanged(false);
          hideModal();
        }, 2000);
      }
      return response.json();
      })
    .then((data) => {
      if (data && data.error) {
        setError(data.error)
        setTimeout(() => {
          setError(null)
        }, 5000);
      }
    })
    .catch((err) => {
      setFetching(false);
      console.log(err);
    })
  }

  const regexCheck = /[\s"]/;
  if (!regexCheck.test(newName) && showRegexWarning) {
    setRegexWarning(false);
  } else if (regexCheck.test(newName) && !showRegexWarning) {
    setRegexWarning(true);
  }
  const disableChange = ((!confirmName || !newName || !newNameConfirm) || showRegexWarning || fetching);


  if (!isVisible) {
    return null;
  }
  return (
    <Modal
      handleClose={hideModal}
      header="Change Password"
    >
      <section className="PasswordModal flex flex-direction--column justify--space-around">
        <div>
          <label>
            Current Password:
          </label>
          <InputText
            inputRef={inputRef}
            label=""
            placeholder="Enter Current Password"
            updateValue={updateConfirmName}
            isPassword
          />
          <label>
            New Password:
          </label>
          <InputText
            inputRef={inputRef}
            label=""
            placeholder="Enter New Password"
            updateValue={updateNewName}
            isPassword
          />
          <label>
            Confirm New Password:
          </label>
          <InputText
            inputRef={inputRef}
            label=""
            placeholder="Confirm New Password"
            updateValue={updateNewNameConfirm}
            isPassword
          />
        </div>
        {
          showRegexWarning && (
            <div className="PasswordModal__warning">
              New password cannot contain whitespace or quotes.
            </div>
          )
        }
        {
          error && (
            <div className="PasswordModal__error">
              {error}
            </div>
          )
        }
        {
          showPasswordMismatch && (
            <div className="PasswordModal__error">
              New passwords does not match. Please try again.
            </div>
          )
        }
        {
          succesfullyChanged && (
            <div className="PasswordModal__success">
              Password changed successfully.
            </div>
          )
        }
        <PrimaryButton
          click={handlePasswordConfirm}
          text="Change Password"
          disabled={disableChange}
        />
      </section>
    </Modal>
  )
}


export default PasswordModal;
