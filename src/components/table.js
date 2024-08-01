import React, { useState } from "react";
import {
  ArrowBackIosRounded,
  ArrowForwardIosRounded,
  RemoveRedEyeRounded,
} from "@mui/icons-material";
import "./style.css";
import CustomButton from "./button";
import InputBox from "./input";
import * as XLSX from "xlsx";


function CustomTable(props) {
  const [pageSize, setPageSize] = useState(10);
  const [currentPage, setCurrentPage] = useState(1);
  const [searchText, setSearchText] = useState("");
  const [filterValue, setFilterValue] = useState("");

  const tableData = props.data;

  const startIndex = (currentPage - 1) * pageSize;
  const endIndex = Math.min(currentPage * pageSize, tableData.length);

  const filteredData = tableData.filter((row) => {
    const searchMatch = Object.values(row).some((value) =>
      String(value).toLowerCase().includes(searchText.toLowerCase())
    );
    const filterMatch =
      filterValue === "" ||
      row[props.filterField].toLowerCase() === filterValue.toLowerCase();

    return searchMatch && filterMatch;
  });

  const visibleData = filteredData.slice(startIndex, endIndex);

  const totalPages = Math.ceil(filteredData.length / pageSize);

  const downloadExcel = () => {
    const worksheet = XLSX.utils.json_to_sheet(filteredData);
    const workbook = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(workbook, worksheet, "Filtered Data");
    XLSX.writeFile(workbook, "mailer.xlsx");
  };
  

  return (
    <div
      className="custom-table-body"
      style={{
        backgroundColor: "white",
        padding: !props.disableHeader ? 20 : 0,
        overflow: "auto",
        borderRadius: 10,
        boxShadow: "rgba(100, 100, 111, 0.1) 0px 7px 29px 0px",
      }}
    >
      <div className="table-header">
        <div
          style={{
            width: 300,
            marginTop: -30,
          }}
          className="search-filter"
        >
          <InputBox
            type="text"
            placeholder="Search..."
            onChange={setSearchText}
          />
        </div>
          <div style={{display:'flex',columnGap:'10px'}}>

        {props.button !== undefined ? (
          <div >
            <CustomButton label={props.buttonLabel} func={props.button} />
            {/* <CustomButton label="Download" func={downloadExcel} /> */}
          </div>
        ) : null}
           
           {props.downloadbutton !== undefined ? (
          <div >
            <CustomButton label={props.downloadbuttonLabel} func={downloadExcel}  />
           
          </div>
            ) : null}
 
          </div>
      </div>
      <div
        style={{
          border: "0.2px solid rgb(216, 216, 216)",
          borderRadius: 8,
          minWidth: "max-content",
          marginTop: 15,
        }}
      >
        <table style={{ border: "none" }}>
          <thead style={{ border: "none" }}>
            <tr style={{ border: "none", backgroundColor: "rgb(248 248 248)" }}>
              <td>S.No</td>
              {props.header.map((item, i) => (
                <td key={i}>{item}</td>
              ))}
              {props.actionType !== "none" ? (
                <td style={{ width: 100 }}>Action</td>
              ) : null}
            </tr>
          </thead>
          <tbody>
            {visibleData.map((row, visibleIndex) => (
              <tr style={{ border: "none" }} key={startIndex + visibleIndex}>
                <td>{startIndex + visibleIndex + 1}</td>
                {props.field.map((item, i) => (
                  <td
                    key={i}
                    dangerouslySetInnerHTML={{ __html: row[item] }}
                  ></td>
                ))}
                {props.actionType !== "none" ? (
                  <td>
                    <div style={{ cursor: "pointer", width: "max-content" }}>
                      {props.actionType === "button" ? (
                         <RemoveRedEyeRounded
                         onClick={() => {
                           if (props.buttonAction !== undefined)
                             props.buttonAction(row);
                         }}
                       />
                      ) : (
                       
                        <CustomButton
                          padding={"8px 15px"}
                          label={props.buttonLabel}
                          onClick={() => {
                            if (props.buttonAction !== undefined)
                              props.buttonAction(row);
                          }}
                        />
                      )}
                    </div>
                  </td>
                ) : null}
              </tr>
            ))}
          </tbody>
        </table>
        <div className="pagination">
          <h4 style={{ fontSize: 14, fontWeight: 500 }}>
            Page {currentPage} of {totalPages}
          </h4>
          <div className="pagination-right">
            <h4 style={{ fontSize: 14, fontWeight: 500 }}>Rows per page:</h4>
            <select
              value={pageSize}
              onChange={(e) => {
                setPageSize(parseInt(e.target.value, 10));
                setCurrentPage(1);
              }}
            >
              <option value={5}>5</option>
              <option value={10}>10</option>
              <option value={20}>20</option>
              <option value={50}>50</option>
            </select>
            <button
              onClick={() => setCurrentPage(currentPage - 1)}
              disabled={currentPage === 1}
            >
              <ArrowBackIosRounded fontSize="sm" />
            </button>
            <button
              onClick={() => setCurrentPage(currentPage + 1)}
              disabled={endIndex >= filteredData.length}
            >
              <ArrowForwardIosRounded fontSize="sm" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default CustomTable;

