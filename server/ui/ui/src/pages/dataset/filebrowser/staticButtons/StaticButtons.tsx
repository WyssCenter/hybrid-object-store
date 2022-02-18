// vendor
import React, { FC, useState } from 'react';
import classNames from 'classnames';
import ReactTooltip from 'react-tooltip';
import { ReactComponent as SearchIcon } from 'Images/icons/icon-search.svg';
import { ReactComponent as FileUpload } from 'Images/icons/icon-file-upload.svg';
import { ReactComponent as FolderAdd } from 'Images/icons/icon-browser-folder-add.svg';
import SearchModal from './searchModal/SearchModal';
// css
import './StaticButtons.scss';


interface Props {
  s3Client: any;
  dataset: string;
  bucket: string;
  refetchS3: any;
  batchUpload:(files: any, prefix: string) => void;
  setNewFolderVisible: (visible: boolean) => void;
  searchFiles: (metadata: any, modifiedBefore: string, modifiedAfter: string) => void;
  searchData: any;
  clearSearchData: () => void;
  loadingSearch: boolean;
  namespace: string;
}


const StaticButtons:FC<Props> = ({
  s3Client,
  dataset,
  bucket,
  refetchS3,
  batchUpload,
  setNewFolderVisible,
  searchFiles,
  searchData,
  clearSearchData,
  loadingSearch,
  namespace,
}: Props) => {
  const [ isModalVisible, setIsModalVisible ] = useState(false);

  const staticButtonsCSS = classNames({
    'StaticButtons flex': true,
    'StaticButtons--search': !!searchData,
  })

  return (
    <div className={staticButtonsCSS}>
      <SearchModal
        hideModal={() => {setIsModalVisible(false)}}
        isVisible={isModalVisible}
        searchFiles={searchFiles}
        dataset={dataset}
        namespace={namespace}
      />
      {
        loadingSearch && (
        <div className="flex flex-1 justify--space-between StaticButtons__searchResults StaticButtons__searchResults--loading">
          <div className="flex StaticButtons__searchResults--loading--text">
            Loading search results...
          </div>
        </div>
        )
      }
      {(searchData !== null) && (
        <div className="flex flex-1 justify--space-between StaticButtons__searchResults">
          <div className="flex">
            {`Viewing Search Results: ${searchData.size} files matching criteria found.`}
          </div>
          <button
            className="StaticButtons__clearSearch"
            onClick={() => clearSearchData()}
          >
            Clear Results
          </button>
        </div>
      )}
      {(searchData === null) && !loadingSearch && (
        <button
          className="StaticButtons__Button StaticButtons__Button--search"
          onClick={() => {setIsModalVisible(true)}}
        >
          <SearchIcon />
          Search
        </button>
      )
      }
      <div
        data-tip="Not supported in search mode."
        role="presentation"
        data-tip-disable={!searchData}
      >
        <button
          className="StaticButtons__Button StaticButtons__Button--add"
          onClick={() => setNewFolderVisible(true)}
          disabled={searchData}
        >
          <FolderAdd />
          New Folder
        </button>
        <ReactTooltip
          place="bottom"
          effect="solid"
        />
      </div>
      <label
        data-tip="Not supported in search mode."
        role="presentation"
        data-tip-disable={!searchData}
        htmlFor="add_file"
        className="flex justify--center"
      >
        <div
          className="StaticButtons__Button StaticButtons__Button--upload"
        >
          <FileUpload />
          Upload Files
        </div>
        <input
        id="add_file"
        className="hidden"
        type="file"
        disabled={searchData}
        onChange={(evt) => {
          const { files } = evt.target;
          batchUpload([files[0]], `${dataset}/`)
        }}
        />
        <ReactTooltip
          place="bottom"
          effect="solid"
        />
      </label>
    </div>
  );
}

export default StaticButtons;
