// vendor
import React,
{
  FC,
  useEffect,
  useContext,
  useReducer,
  useRef,
} from 'react';
// components
import Dropdown from 'Components/dropdown/index';
// context
import AppContext from 'Src/AppContext';
// css
import './PermissionDropdown.scss';


interface DropdownItem {
  name: string;
}

interface Props {
  name: string;
  permission: DropdownItem;
  updatePermissions: (permission: DropdownItem, name: string, errorCallback: () => void) => void;
}


interface State {
  selectedItem?: DropdownItem;
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
        ...state
      }
  }
}

const PermissionDropdown: FC<Props> = ({
  name,
  permission,
  updatePermissions
}: Props) => {
  // context
  const { user } = useContext(AppContext);
  // state
  const [state, dispatch] = useReducer(
    stateReducer,
    {
      dropdownVisibility: false,
      selectedItem: permission,
    });
  // refs
  const dropdownRef = useRef(null);

  useEffect(() => {
    if (permission.name !== state.selectedItem.name) {
      updatePermissions(state.selectedItem, name, () => {
        dispatch({
          type: 'set',
          payload: { selectedItem: permission }
        });
      });
    }
  }, [state.selectedItem, updatePermissions, name, permission]);

  const authorized = user.profile.role !== 'user';

  let permName = 'Read Only';
  if (state.selectedItem.name === 'rw' || state.selectedItem.name === 'Read & Write') {
    permName = "Read & Write";
  }

  return (
    <div
      className="PermissionDropdown"
      ref={dropdownRef}
    >
      {
        authorized && (
        <Dropdown
          customStyle="small"
          dispatch={dispatch}
          label={permName}
          listItems={[{name: 'Read & Write'}, {name: 'Read Only'}]}
          visibility={state.dropdownVisibility}
        />
        )
      }
      { !authorized && permName }
    </div>
  )

}


export default PermissionDropdown;
