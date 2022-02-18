// vendor
import React,
{
  FC,
  useState,
  useCallback,
  useContext,
} from 'react';
import { Link, useParams } from 'react-router-dom';
import { faChevronRight, faChevronDown } from '@fortawesome/free-solid-svg-icons';
import moment from 'moment';
// context
import AppContext from 'Src/AppContext';
// componets
import { IconButton } from 'Components/button/index';
import DatasetDetails from './details/DatasetDetails';
import FileBrowserButton from './buttons/FileBrowserButton';
// context
import NamespaceContext from '../../NamespaceContext';
// css
import { ReactComponent as Disabled } from 'Images/icons/icon-sync-disabled.svg';
import { ReactComponent as Simplex } from 'Images/icons/icon-sync-simplex.svg';
import { ReactComponent as Duplex } from 'Images/icons/icon-sync-duplex.svg';
import './DatasetItem.scss';


type Data = {
  namespace?: string;
  error?: string;
}


interface Dataset {
  created: string;
  description: string;
  directory: string;
  name: string;
  sync: string;
  showDeleteBadge: boolean;
}

interface ObjectStore {
  type: string;
}

interface Namespace {
  bucket_name: string;
  name: string;
  object_store: ObjectStore;
}

interface Props {
  dataset: Dataset;
  namespaceData: Namespace
}

interface ParamTypes {
  namespace: string;
}

const DatsetItem: FC<Props> = ({ dataset, namespaceData }: Props ) => {
  // context
  const { user } = useContext(AppContext);
  // state
  const [isExpanded, setIsExpanded] = useState(false);
  // params
  const { namespace } = useParams<ParamTypes>();

  /**
  * Method toggles detail view
  * @param {}
  * @return {void}
  * @calls {state#setIsExpanded}
  */
  const toggleView = useCallback(() => {
    setIsExpanded(!isExpanded);
  }, [isExpanded])

  const icon = isExpanded ? faChevronDown : faChevronRight;

  const permissions = user.profile.role !== 'user';

  let sanitizedSyncName = 'Disabled';
  if (dataset.sync === 'simplex') {
    sanitizedSyncName = '1-Way';
  } else if (dataset.sync === 'duplex') {
    sanitizedSyncName = '2-Way';
  }

  return (
    <>
      <tr>
        <td className="DatasetItem__td--empty">
          <IconButton
            click={toggleView}
            color="primary"
            backgroundColor="transparent"
            icon={icon}
          />
        </td>
        <td>
          <Link to={`/${namespace}/${dataset.name}`}>
            {dataset.name}
          </Link>
          {
            dataset.showDeleteBadge && (
              <span className="DatasetItem__badge">
                DELETE SCHEDULED
              </span>
            )
          }
        </td>
        <td>{dataset.description}</td>
        <td>{moment.utc(dataset.created).fromNow()}</td>
        <td className={`DatasetItem__sync`}>
          {(dataset.sync === 'disabled') && <Disabled />}
          {(dataset.sync === 'simplex') && <Simplex />}
          {(dataset.sync === 'duplex') && <Duplex />}
          {sanitizedSyncName}
        </td>
      </tr>

      <DatasetDetails
        datasetName={dataset.name}
        isExpanded={isExpanded}
      />

    </>
  );
}


export default DatsetItem;
