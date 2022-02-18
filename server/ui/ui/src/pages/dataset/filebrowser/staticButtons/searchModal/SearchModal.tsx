// vendor
import React, { FC, useState, useRef, useEffect } from 'react';
import Modal from 'Components/modal/Modal';
import { get } from 'Environment/createEnvironment';
// components
import InputText from 'Components/form/text/index';
import { ReactComponent as Add } from 'Images/icons/add-button.svg';
import { PrimaryButton, TernaryButton } from 'Components/button/index';

// css
import './SearchModal.scss';


interface Props {
  hideModal: () => void;
  isVisible: boolean;
  searchFiles: (metadata: any, modifiedBefore: string, modifiedAfter: string) => void;
  namespace: string;
  dataset: string;
}


const SearchModal:FC<Props> = ({
    hideModal,
    isVisible,
    searchFiles,
    namespace,
    dataset,
}: Props) => {


  const [newMetaKey, setNewMetaKey] = useState('');
  const [newMetaValue, setNewMetaValue] = useState('');
  const [newBeforeModified, setBeforeModified] = useState(null);
  const [newAfterModified, setAfterModified] = useState(null);
  const [metaObject, setMetaObject] = useState(null);
  const [autoCompleteKeys, setAutoCompleteKeys] = useState([])
  const [autoCompleteValues, setAutoCompleteValues] = useState([])
  const [showKeysAuto, setShowKeysAuto] = useState(false);
  const [showValuesAuto, setShowValuesAuto] = useState(false);
  const [hasError, setHasError] = useState(false);

  const keyRef = useRef(null);
  const valueRef = useRef(null);
  const beforeRef = useRef(null);
  const afterRef = useRef(null);

  useEffect(() => {
    if(!isVisible) {
      setMetaObject(null);
      setNewMetaValue(null);
      setNewMetaKey(null);
      setBeforeModified(null);
      setAfterModified(null);
      if(keyRef.current) {
        keyRef.current.value = '';
      }
      if(afterRef.current) {
        valueRef.current.value = '';
      }
    }
  }, [isVisible])

  useEffect(() => {
    if(newMetaKey) {
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

  if(!isVisible){
    return null;
  }



  const handleSearch = () => {
    if (
      (!newBeforeModified && newAfterModified)
      || (newBeforeModified && !newAfterModified)
      || (newBeforeModified >= newAfterModified)
      ) {
      searchFiles(metaObject, newBeforeModified, newAfterModified);

      hideModal();
    } else {
      setHasError(true);
      setTimeout(() => {
        setHasError(false);
      }, 5000);
    }
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
  const updateBeforedate = (evt: any) => {
    setBeforeModified(evt.target.value);
  }
  const updateAfterdate = (evt: any) => {
    setAfterModified(evt.target.value);
  }

  const enableSearch = ((metaObject && (Object.keys(metaObject).length > 0))
    || newAfterModified
    || newBeforeModified)


  return (
    <Modal
      handleClose={hideModal}
      header="Dataset Search"
      size="large-full"
      noCancel
    >
      <section className="SearchModal">
        <div className="SearchModal__primary">
          Search for files in this dataset based on metadata key-value pairs and/or modified date.
        </div>
        <div className="SearchModal__secondary">
          If no metadata values are set, all data matching the modified date values will be returned.
        </div>
        <div className="SearchModal__secondary">
          If “modified before” and “modified after” are not set, any file that contains the specified key-value pairs will be returned.
        </div>
        <div className="SearchModal__inputs flex justify--space-around">
          <div className="flex align-items--center">
            <span>
              Metadata Key:
            </span>
            <div className="SearchModal__input">
              <InputText
                inputRef={keyRef}
                label=""
                placeholder=""
                updateValue={updateMetaKey}
                onFocus={() => setShowKeysAuto(true)}
                onBlur={(evt: any) => {
                  if (!evt.relatedTarget || (evt.relatedTarget.className !== 'SearchModal__auto--button')) {
                  setShowKeysAuto(false)
                  }
                }
                }
              />
              {
                showKeysAuto && autoCompleteKeys && (autoCompleteKeys.length > 0) && (
                <div className="SearchModal__auto">
                  {
                    autoCompleteKeys.map((key) => {
                      return (
                        <button
                          className="SearchModal__auto--button"
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
            <div className="SearchModal__input">
              <InputText
                inputRef={valueRef}
                label=""
                placeholder=""
                updateValue={updateMetaValue}
                onFocus={() => setShowValuesAuto(true)}
                onBlur={(evt: any) => {
                  if (!evt.relatedTarget || (evt.relatedTarget.className !== 'SearchModal__auto--button')) {
                  setShowValuesAuto(false)
                  }
                }
              }
              />
              {
                showValuesAuto && autoCompleteValues && (autoCompleteValues.length > 0) && (
                <div className="SearchModal__auto">
                  {
                    autoCompleteValues.map((key) => {
                      return (
                        <button
                          className="SearchModal__auto--button"
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
              className="SearchModal__add"
              onClick={() => handleAdd()}
            >
              <Add />
            </button>
          </div>
        </div>
        {
          !metaObject && (
            <div className="SearchModal__placeholder">
              <div className="SearchModal__placeholder--primary">
                No Metadata Parameters Set
              </div>
              <div className="SearchModal__placeholder--secondary">
                Add key-value pairs above to search for matching files
              </div>
            </div>
          )
        }
        {
          metaObject && (Object.keys(metaObject).length > 0)
          && (
              <div className="SearchModal__table">
                <div className="SearchModal__header">
                  <span>
                    Key
                  </span>
                  <span>
                    Value
                  </span>
                </div>
                <div className="SearchModal__contents">
                  {
                    metaObject && Object.keys(metaObject).map((entry) => (
                    <div className="SearchModal__entry" key={entry}>
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
        <div className="SearchModal__inputs flex justify--space-around">
          <div className="flex align-items--center">
            <span>
              Modified Before:
            </span>
            <input
              type="datetime-local"
              ref={beforeRef}
              onChange={updateBeforedate}
            />
          </div>
          <div className="flex align-items--center">
            <span>
              Modified After:
            </span>
            <input
              type="datetime-local"
              ref={afterRef}
              onChange={updateAfterdate}
            />
          </div>
        </div>
        {
          hasError && (
            <div className="SearchModal__error">
              Please select a valid time range and try again.
            </div>
          )
        }
        <div className="SearchModal__buttons flex justify--flex-end">
          <TernaryButton
            click={hideModal}
            text="Cancel"
          />
          <PrimaryButton
            disabled={!enableSearch}
            click={handleSearch}
            text="Search"
          />
        </div>
      </section>
    </Modal>
  );
}

export default SearchModal;
