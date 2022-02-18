// vendor
import React, {
    FC,
    KeyboardEventHandler,
    useCallback,
    useRef,
  } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSearch } from '@fortawesome/free-solid-svg-icons';
// css
import './Search.scss';


// needs to be generic so we can reuse
interface ArrayObject {
  [key: string]: any;
}


interface Props {
  disabled: boolean,
  placeholder: string,
  list: Array<ArrayObject>,
  updateList: (list: any) => void,
}

const Search: FC<Props> = ({
  disabled,
  placeholder,
  list,
  updateList,
}:Props) => {
  const inputRef = useRef(null);

  /**
  * Method updates the search value;
  * @param {KeyboardEvent} event
  * @return {void}
  * @fires {#setSearch}
  */
  const handleSearchEvent:any = useCallback((event: KeyboardEvent) => {
    const element = event.currentTarget as HTMLInputElement
    const value = element.value


    const newList = list.filter(entry => {
      const stringifiedEntry = JSON.stringify(Object.values(entry));
      if (stringifiedEntry.indexOf(value) > -1) {
        return true;
      }
      return false;
    });

    updateList(newList);

    return;
  }, [list, updateList])

  return (
    <div className="Search column-3-span-4">
      <label className="InputText__label">Search</label>
      <input
        className="InputText__input"
        disabled={disabled}
        onKeyUp={(event) => handleSearchEvent(event)}
        placeholder={placeholder || ''}
        ref={inputRef}
        type="text"
      >

      </input>
      <div className="Search__icon"> 
        <FontAwesomeIcon icon={faSearch} color="#0b1425" />
      </div>
    </div>
  )
}


export default Search;
