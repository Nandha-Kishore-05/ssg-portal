import Select from "react-select";
import "./style.css";

function CustomSelect(props) {
  return (
    <div
      className="input-box"
      style={{
        marginTop: props.margin ? props.margin :null,
        flex: !props.widthFull ? 3 : null,
      }}
    >
      {props.label ? <label>{props.label}</label> : null}
      <div style={{ marginTop: 5, fontSize: 16 ,fontWeight:"bolder"}}>
        <Select
          onChange={(e) => {
            props.isMulti ? props.onChange(e) : props.onChange(e);

          }}
          options={props.options}
          isSearchable={true}
          placeholder={props.placeholder}
          isMulti={props.isMulti ? true : false}
          value={props.value}
        />
      </div>
    </div>
  );
}

export default CustomSelect;
