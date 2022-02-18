// vendor
import React,
{
  FC,
  useCallback,
  useEffect,
  useState,
} from 'react';
import {
  useParams,
} from "react-router-dom"
// environment
import { get } from 'Environment/createEnvironment';
// css
import  './DatasetDetails.scss';

interface Group {
  group_name: string;
}

interface Permission {
  group: Group;
  permission: 'r' | 'rw';
}


interface Data {
  namespace?: string;
  error?: string;
  permissions?: Array<Permission>;
}


interface Props {
  datasetName: string;
  isExpanded: boolean;
}


interface ParamTypes {
  namespace: string;
}



const DatasetExpandedItem: FC<Props> = ({ datasetName, isExpanded }: Props ) => {
  // dataset
  const [datasetData, setDatasetData] = useState(null);
  // params
  const { namespace } = useParams<ParamTypes>();

  /**
  * Method gets dataset data for rendering
  * @param {}
  * @return {void}
  * @calls {environment#get}
  * @calls {state#setDatasetData}
  */
  const fetchDatasetData = useCallback(() => {
    if (isExpanded && (datasetData == null)){
      get(`namespace/${namespace}/dataset/${datasetName}`).then((response: Response) => {
          return response.json();
        })
        .then((data: Data) => {
          setDatasetData(data);

        })
    }
  }, [namespace, datasetName, datasetData, isExpanded])

  useEffect(() => {
    fetchDatasetData();
  }, [fetchDatasetData, isExpanded]);

  if (!isExpanded) {
    return (null);
  }

  const rwList = datasetData && datasetData.permissions.filter((permission: Permission) => {
    return permission.permission === 'rw';
  }).map((permission: Permission) => {
    return permission.group.group_name.replace('-hoss-default-group', '')
  }).sort().join(', ');

  const rList = datasetData && datasetData.permissions.filter((permission: Permission) => {
    return permission.permission === 'r';
  }).map((permission: Permission) => {
    return permission.group.group_name.replace('-hoss-default-group', '')
  }).sort().join(', ');



  return (
    <tr>
      <td
        className="DatasetDetails"
        colSpan={6}
      >
        { datasetData &&
          <div className="DatasetDetails__permission grid">
            <p className="column-1-span-2">
              <b>Owner:</b>
              {` ${datasetData.owner.username}`}
            </p>

            { rwList  &&
              <p className="column-3-span-4">
                <b>rw:</b>
                {` ${rwList}`}
              </p>
            }

            { rList &&
              <p className="column-3-span-4">
                <b>r:</b>
                {` ${rList}`}
              </p>
            }

          </div>
        }

        { !datasetData && (

          <p>Loading</p>

        )}

      </td>
    </tr>
  );
}


export default DatasetExpandedItem;
