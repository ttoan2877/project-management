import React, { useState, useMemo } from 'react'
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  TouchableOpacity,
} from 'react-native'
import SearchInput from './SearchInput'
import { get } from 'lodash'

const ListEmptyComponent = ({ emptyText }) => (
  <View style={styles.emptyContainer}>
    <Text style={styles.text}>{emptyText}</Text>
  </View>
)

const CommonLayout = ({
  data,
  renderItem,
  emptyText,
  onRefresh,
  refreshing,
  onCreate,
  filterValue,
}) => {
  const [searchText, setSearchText] = useState('')

  const listData = useMemo(
    () =>
      searchText.length > 0
        ? data.filter((x) =>
            get(x, filterValue)
              .toLowerCase()
              .includes(searchText.toLowerCase()),
          )
        : data,
    [searchText, data, filterValue],
  )
  return (
    <View style={styles.container}>
      <SearchInput searchText={searchText} setSearchText={setSearchText} />
      {listData.length > 0 ? (
        <FlatList
          refreshing={refreshing}
          onRefresh={onRefresh}
          style={styles.list}
          data={listData}
          renderItem={renderItem}
        />
      ) : (
        <ListEmptyComponent emptyText={emptyText} />
      )}
      {onCreate && (
        <TouchableOpacity onPress={onCreate} style={styles.addBtn}>
          <Text style={styles.btnLabel}>ADD</Text>
        </TouchableOpacity>
      )}
    </View>
  )
}

export default CommonLayout

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: 'black',
  },
  list: {
    flex: 1,
    marginBottom: 32,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  addBtn: {
    width: 48,
    height: 48,
    borderRadius: 48,
    backgroundColor: 'firebrick',
    position: 'absolute',
    bottom: 60,
    right: 16,
    justifyContent: 'center',
    alignItems: 'center',
    borderColor: 'gray',
    borderWidth: 4,
  },
  btnLabel: {
    color: 'azure',
    fontWeight: 'bold',
  },
  text: {
    fontWeight: 'bold',
    color: 'azure',
  },
})
