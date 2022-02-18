// vendor
import React,
{
  useEffect,
  useReducer,
  useRef,
  FC,
  KeyboardEvent,
} from 'react';
// components
import Dropdown from 'Components/dropdown/index';
import InputText from 'Components/form/text/index';
// css
import './NamespaceFields.scss';

interface ObjectStore {
  name: string;
  type: string;
}


interface NamespaceFieldsProps {
  handleBucketNameEvent: (event: KeyboardEvent) => void;
  handleObjectChangeEvent: (value: ObjectStore) => void;
  objectStoreList: Array<ObjectStore>;
}


interface State {
  selectedItem?: ObjectStore;
  dropdownVisibility?: boolean;
}

interface Action {
  type: string;
  payload?: State;
}


const stateReducer = (state: State, action: Action) => {
  switch(action.type) {
    case 'toggle':
      return {
        dropdownVisibility: !state.dropdownVisibility,
        selectedItem: state.selectedItem,
      };
    case 'set':
      return {
        dropdownVisibility: false,
        ...action.payload
      };
    case 'close':
      return {
        dropdownVisibility: false,
        selectedItem: state.selectedItem,
      };
    default:
      return {
        ...state,
      }
  }
}

const NamespaceFields: FC<NamespaceFieldsProps> = ({
  handleBucketNameEvent,
  handleObjectChangeEvent,
  objectStoreList,
}: NamespaceFieldsProps) => {

  //refs
  const bucketNameRef = useRef(null);
  // reducers
  const [dropdownState, dispatch] = useReducer(
    stateReducer,
    {
      dropdownVisibility: false,
      selectedItem: objectStoreList[0],
    });

    useEffect(()=> {
      handleObjectChangeEvent(dropdownState.selectedItem);
    }, [dropdownState.selectedItem, handleObjectChangeEvent]);

  const label = `${dropdownState.selectedItem.name} (${dropdownState.selectedItem.type})`;

  return (
    <>
      <InputText
        css=""
        isRequired
        inputRef={bucketNameRef}
        label="Bucket Name"
        updateValue={handleBucketNameEvent}
      />
      <div className="NamespaceFields__dropdown-container">
        <p className="NamespaceFields__p">Object Store Type</p>
        <Dropdown
          customStyle="menu-right"
          dispatch={dispatch}
          label={label}
          listItems={objectStoreList}
          visibility={dropdownState.dropdownVisibility}
        />
      </div>
    </>
  )
}

export default NamespaceFields;
