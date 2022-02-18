// vendor
import React,
{
  FC,
  useCallback,
  useEffect,
  useState,
} from 'react';
// environment
import { get } from 'Environment/createEnvironment';
// css
import  './NamespaceDetails.scss';

interface Dataset {
  name: string;
}


type Data = Array<Dataset>;

interface Props {
  namespaceName: string;
  isExpanded: boolean;
}



const NamespaceDetails: FC<Props> = ({ namespaceName, isExpanded }: Props ) => {
  // dataset
  const [namespaceData, setNamespaceData] = useState(null);

  /**
  * Method gets namespace data for rendering
  * @param {}
  * @return {void}
  * @calls {environment#get}
  * @calls {state#setDatasetData}
  */
  const fetchDatasetData = useCallback(() => {
    if (isExpanded && (namespaceData == null)){
      get(`namespace/${namespaceName}/dataset/`).then((response: Response) => {
          return response.json();
        })
        .then((data: Data) => {
          setNamespaceData(data);
        });
    }
  }, [namespaceName, namespaceData, isExpanded])

  useEffect(() => {
    fetchDatasetData();
  }, [fetchDatasetData, isExpanded]);

  if (!isExpanded) {
    return (null);
  }

  const datasets = namespaceData && namespaceData.map((dataset: Dataset) => {
    return dataset.name;
  }).sort().join(', ');


  const datasetText = datasets && datasets.length > 0
    ? datasets
    : 'There are no datasets in this namespace.'


  return (
    <tr>
      <td
        className="NamespaceDetails"
        colSpan={7}
      >
        {
          <p>
            <b>Datasets:</b>
              {` ${datasetText}`}
          </p>
        }

      </td>
    </tr>
  );
}


export default NamespaceDetails;
