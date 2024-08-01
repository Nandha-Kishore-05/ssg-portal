import "./style.css";
import Select from "react-select";

function SelectBox(props) {
  return (
    <div className="input-box" style={{ marginTop: props.margin }}>
      <label>{props.label}</label>
      <br />
      <Select
        styles={{
          control: (baseStyles) => ({
            ...baseStyles,
            backgroundColor: "#F7F7F8",
            border: "none",
            padding: 5,
            marginTop: 8,
            borderRadius: 8,
            outline:"none",
          }),
        }}
        onChange={(e) => {
          props.onChange(e.value);
        }}
        options={props.options}
        isSearchable={true}
        placeholder={props.placeholder}
      />
    </div>
  );
}

export default SelectBox;
