// vendor
import React, { FC, useState, useRef, useEffect } from 'react';
import Modal from 'Components/modal/Modal';
import { get } from 'Environment/createEnvironment';
// components
import InputText from 'Components/form/text/index';
import { ReactComponent as Add } from 'Images/icons/add-button.svg';
import { PrimaryButton, TernaryButton } from 'Components/button/index';

// css
import './TagsModal.scss';


interface Props {
  hideModal: () => void;
  isVisible: boolean;
  setPendingUpload: (data: any) => void;
  upload: (files: any, path: string, metaObject: any) => void;
  pendingUpload: any;
  setIsLocked: (locked: boolean) => void;
  namespace: string;
  dataset: string;
}


const TagsModal:FC<Props> = ({
    hideModal,
    isVisible,
    setIsLocked,
    upload,
    pendingUpload,
    setPendingUpload,
    namespace,
    dataset,
}: Props) => {


  const [newMetaKey, setNewMetaKey] = useState('');
  const [newMetaValue, setNewMetaValue] = useState('');
  const [metaObject, setMetaObject] = useState({});
  const [autoCompleteKeys, setAutoCompleteKeys] = useState([])
  const [autoCompleteValues, setAutoCompleteValues] = useState([])
  const [showKeysAuto, setShowKeysAuto] = useState(false);
  const [showValuesAuto, setShowValuesAuto] = useState(false);

  const keyRef = useRef(null);
  const valueRef = useRef(null);


  useEffect(() => {
    if(newMetaKey){
      get(`search/namespace/${namespace}/dataset/${dataset}/key?prefix=${newMetaKey}&limit=10`)
      .then((res: any) => res.json())
      .then((data: any) => {
        if (!data || !data.keys) {
          setAutoCompleteKeys([])
        } else {
          setAutoCompleteKeys(data.keys);
        }
      });
    } else {
      setAutoCompleteKeys([]);
    }
  }, [newMetaKey])

  useEffect(() => {
    if(newMetaKey) {
      get(`search/namespace/${namespace}/dataset/${dataset}/key/${newMetaKey}/value?prefix=${newMetaValue}&limit=10`)
      .then((res: any) => res.json())
      .then((data: any) => {
        if (!data || !data.values) {
          setAutoCompleteValues([])
        } else {
          setAutoCompleteValues(data.values);
        }
      });
    }
  }, [newMetaValue])

  useEffect(() => {
    if(!isVisible) {
      setMetaObject(null);
      setNewMetaValue('');
      setNewMetaKey('');
    }
  }, [isVisible])

  if(!isVisible){
    return null;
  }

  const handleAdd = () => {
    keyRef.current.value = '';
    valueRef.current.value = '';
    setMetaObject((prevState: any) => ({
      ...prevState,
      [newMetaKey]: newMetaValue,
    }))
    setNewMetaKey('');
    setNewMetaValue('');
  }

  const updateMetaKey = (evt: any) => {
    if((evt.key === 'Enter') && newMetaKey && (newMetaKey.length > 0) && newMetaValue && (newMetaValue.length > 0)){
      handleAdd();
    } else {
      setNewMetaKey(evt.target.value);
    }
  }

  const updateMetaValue = (evt: any) => {
    if((evt.key === 'Enter') && newMetaKey && (newMetaKey.length > 0) && newMetaValue && (newMetaValue.length > 0)){
      handleAdd();
    } else {
      setNewMetaValue(evt.target.value);
    }
  }

  return (
    <Modal
      handleClose={() => {
        setIsLocked(false);
        setPendingUpload(null);
        hideModal();
      }}
      header="Add File Metadata"
      size="large-full"
    >
      <section className="TagsModal">
        <div className="TagsModal__primary">
          (Optional) Provide key-value pairs to be set as metadata on all uploaded files.
        </div>
        <div className="TagsModal__secondary">
          Metadata can be used to store additional information related to file and is useful when searching for data.
        </div>
        <div className="TagsModal__secondary">
          Adding metadata is not required for uploading.
        </div>
        <div className="TagsModal__inputs flex justify--space-around">
          <div className="flex align-items--center">
            <span>
              Metadata Key:
            </span>
            <div className="TagsModal__input">
              <InputText
                inputRef={keyRef}
                label=""
                placeholder=""
                updateValue={updateMetaKey}
                onFocus={() => setShowKeysAuto(true)}
                onBlur={(evt: any) => {
                  if (!evt.relatedTarget || (evt.relatedTarget.className !== 'TagsModal__auto--button')) {
                  setShowKeysAuto(false)
                  }
                }
                }
              />
              {
                showKeysAuto && autoCompleteKeys && (autoCompleteKeys.length > 0) && (
                <div className="TagsModal__auto">
                  {
                    autoCompleteKeys.map((key) => {
                      return (
                        <button
                          className="TagsModal__auto--button"
                          key={key}
                          onClick={() => {
                            setNewMetaKey(key);
                            keyRef.current.value = key;
                            setShowKeysAuto(false);
                          }}
                        >
                          {key}
                        </button>
                      )
                    })
                  }
                </div>
                )
              }
            </div>
          </div>
          <div className="flex align-items--center">
            <span>
              Metadata Value:
            </span>
            <div className="TagsModal__input">
              <InputText
                inputRef={valueRef}
                label=""
                placeholder=""
                updateValue={updateMetaValue}
                onFocus={() => setShowValuesAuto(true)}
                onBlur={(evt: any) => {
                  if (!evt.relatedTarget || (evt.relatedTarget.className !== 'TagsModal__auto--button')) {
                  setShowValuesAuto(false)
                  }
                }
              }
              />
              {
                showValuesAuto && autoCompleteValues && (autoCompleteValues.length > 0) && (
                <div className="TagsModal__auto">
                  {
                    autoCompleteValues.map((key) => {
                      return (
                        <button
                          className="TagsModal__auto--button"
                          key={key}
                          onClick={() => {
                            setNewMetaValue(key);
                            valueRef.current.value = key;
                            setShowValuesAuto(false);
                          }}
                        >
                          {key}
                        </button>
                      )
                    })
                  }
                </div>
                )
              }
            </div>
            <button
              className="TagsModal__add"
              onClick={() => handleAdd()}
            >
              <Add />
            </button>
          </div>
        </div>
        {
          !metaObject && (
            <div className="TagsModal__placeholder">
              <div className="TagsModal__placeholder--primary">
                No Metadata Provided
              </div>
              <div className="TagsModal__placeholder--secondary">
                Adding metadata is not required to upload.
              </div>
            </div>
          )
        }
        {
          metaObject && (Object.keys(metaObject).length > 0)
          && (
              <div className="TagsModal__table">
                <div className="TagsModal__header">
                  <span>
                    Key
                  </span>
                  <span>
                    Value
                  </span>
                </div>
                <div className="TagsModal__contents">
                  {
                    metaObject && Object.keys(metaObject).map((entry) => (
                    <div className="TagsModal__entry" key={entry}>
                      <span>
                        {entry}
                      </span>
                      <span>
                        {metaObject[entry]}
                      </span>
                    </div>
                    ))
                  }
                </div>
              </div>
          )
        }
        <div className="TagsModal__buttons flex justify--flex-end">
          <TernaryButton
            click={() => {
              setIsLocked(false);
              setPendingUpload(null);
              hideModal();
            }}
            text="Cancel"
          />
          <PrimaryButton
            click={() => {
              upload(
                pendingUpload.files,
                pendingUpload.path,
                metaObject,
              )
              hideModal();
            }}
            text="Upload"
          />
        </div>
      </section>
    </Modal>
  );
}

export default TagsModal;
