// vendor
import React, { FC, useState, useRef } from 'react';
import Moment from 'moment';
import ReactTooltip from 'react-tooltip';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { InputText } from 'Components/form/text/index';
import { faCopy } from '@fortawesome/free-solid-svg-icons';
import { ReactComponent as CloseIcon } from 'Images/icons/icon-close.svg';


// css
import './Details.scss';


interface Props {
  data: any;
  setFocus: any;
  focusedFile: any;
  humanFileSize: any;
  setFocusedFileData: (data: any) => void;
  dataset: string;
  namespace: string;
  isVisible: boolean;
  setDetailsVisible: (visible: boolean) => void;
  detailsVisible: boolean;
}


const Details:FC<Props> = ({
  data,
  setFocus,
  setFocusedFileData,
  focusedFile,
  humanFileSize,
  dataset,
  namespace,
  isVisible,
  setDetailsVisible,
  detailsVisible,
}: Props) => {

  const [eTagCopied, setETagCopied] = useState(false);
  const [URICopied, setURICopied] = useState(false);
  const [copiedEntry, setCopiedEntry] = useState(null);
  const [copiedValue, setCopiedValue] = useState(null);
  const [metaSearch, setMetaSearch] = useState('');
  const inputRef = useRef(null);

  if (data && !data.ETag) {
    return null;
  }

  if (!isVisible && !detailsVisible) {
    return (
      <div
        className="Details__collapsed"
        onClick={() => {
          setDetailsVisible(true);
        }}
      >
      </div>
    )
  }
  if (!isVisible && detailsVisible) {
    return (
      <div
        className="Details__expanded"
      >
      <div
        className="Details__close"
        onClick={() => {
          setDetailsVisible(false);
        }}
      >
        <CloseIcon />
      </div>
      <div className="Details__header">
        File Details
      </div>
      <div className="Details__filler">
        Select a file in the browser to view attributes and metadata
      </div>
      </div>
    )
  }

  const root = document.documentElement;
  const primaryHash = root.style.cssText ?
    root.style.cssText.split(';')[0].split(':')[1]
    : '#957299';

  const formattedETag = data.ETag.substring(1, data.ETag.length -1);

  const copyToClipboard = (text: string) => {
    const temp = document.createElement('textarea');
    document.body.appendChild(temp);
    temp.value = text;
    temp.select();
    document.execCommand("copy");
    document.body.removeChild(temp);
  };

  const name = focusedFile.Key.split('/').slice(1).join('/');

  const getURI = () => {
    return `hoss+${window.location.origin}:${namespace}:${focusedFile.Key}`
  };

  const hossURI = URICopied
    ? 'Copied to clipboard!'
    : `hoss+${window.location.origin}:${namespace}:${focusedFile.Key}`;

  const eTAG = eTagCopied
    ? 'Copied to clipboard!'
    : formattedETag;

  const filterMeta = (tags: any) => {
    if(metaSearch === '') {
      return tags;
    } else {
      const newTagObject = {};
      Object.keys(tags).forEach((entry: string) => {
        if  (
          ((entry.toLowerCase().indexOf(metaSearch.toLowerCase())) > -1)
          || (tags[entry].toString().indexOf(metaSearch.toLowerCase()) > -1)
        ) {
          newTagObject[entry] = tags[entry];
        }
      })
      return newTagObject;
    }
  };


  const filteredMeta = filterMeta(data.Metadata);


  return (
    <div className="Details">
      <div
        className="Details__close"
        onClick={() => {
          setFocus(null);
          setFocusedFileData(null);
        }}
      >
        <CloseIcon />
      </div>
      <div className="Details__header">
        File Details
      </div>
      <div className="Details__container">
        <div className="Details__title">Key</div>
        <div className="Details__data">{name}</div>
        <div className="Details__title">Last Modified</div>
        <div className="Details__data">{Moment(focusedFile.LastModified).format('MMM. D, YYYY [at] h:mm A z')}</div>
        <div className="Details__title">Size</div>
        <div className="Details__data">{humanFileSize(focusedFile.Size)}</div>
        <div className="Details__title">Entity tag (ETag)</div>
        <div className="Details__data Details__data--small">
          <span onClick={() => {
            copyToClipboard(formattedETag);
            setETagCopied(true);
            setTimeout(() => {
              setETagCopied(false);
            }, 2000)
          }}>
            <FontAwesomeIcon icon={faCopy} color={primaryHash} />
          </span>
          {eTAG}
        </div>
        <div className="Details__title">Hoss URI</div>
        <div className="Details__data Details__data--small">
          <span onClick={() => {
            copyToClipboard(hossURI);
            setURICopied(true);
            setTimeout(() => {
              setURICopied(false);
            }, 2000)
          }}>
            <FontAwesomeIcon icon={faCopy} color={primaryHash} />
          </span>
          {hossURI}
        </div>
        <div className="Details__meta">
          <div className="Details__meta__search flex">
            <div className="Details__title">Metadata</div>
              {(Object.keys(data.Metadata).length > 0) && (
                <InputText
                  css="small"
                  label=""
                  inputRef={inputRef}
                  placeholder=""
                  updateValue={(evt: Event) => setMetaSearch((evt.target as HTMLInputElement).value)}
                />
              )}
          </div>
          {
            (Object.keys(filteredMeta).length === 0) && Object.keys(data.Metadata).length === 0
            && (
              <div className="Details__data--small">This file has no metadata tags</div>
            )
          }
          {
            (Object.keys(filteredMeta).length === 0) && Object.keys(data.Metadata).length > 0
            && (
              <div className="Details__data--small">No matching values</div>
            )
          }
          {
            (Object.keys(filteredMeta).length > 0)
            && (
              <div className="Details__meta__table">
                <div className="Details__meta__header">
                  <span>
                    Key
                  </span>
                  <span>
                    Value
                  </span>
                </div>
                <div className="Details__meta__contents">
                  {
                    Object.keys(filteredMeta).map((entry, index) => (
                    <div className="Details__meta__entry" key={entry}>
                      <span
                        key={copiedEntry === index ? `${entry}--hidden` : entry}
                        data-tip={copiedEntry === index ? 'Copied to clipboard' : entry.length > 15 ? entry : ''}
                        role="presentation"
                        onClick={() =>{
                          copyToClipboard(entry);
                          setCopiedEntry(index);
                          setTimeout(() => {
                            setCopiedEntry(null)
                          }, 3000);
                        }}
                      >
                        {entry}
                       <FontAwesomeIcon icon={faCopy} color={primaryHash} />

                        <ReactTooltip
                          place="bottom"
                          effect="solid"
                        />
                      </span>
                      <span
                        key={copiedValue === index ? `${data.Metadata[entry]}--hidden` : data.Metadata[entry]}
                        role="presentation"
                        onClick={() =>{
                          copyToClipboard(data.Metadata[entry])
                          setCopiedValue(index);
                          setTimeout(() => {
                            setCopiedValue(null)
                          }, 3000);
                        }}
                        data-tip={copiedValue === index ? 'Copied to clipboard' : data.Metadata[entry].length > 15 ? data.Metadata[entry] : ''}
                      >
                        {data.Metadata[entry]}
                       <FontAwesomeIcon icon={faCopy} color={primaryHash} />

                        <ReactTooltip
                          place="bottom"
                          effect="solid"
                        />
                      </span>
                    </div>
                    ))
                  }
                </div>
              </div>
            )
          }
        </div>
      </div>
    </div>
  );
}

export default Details;
