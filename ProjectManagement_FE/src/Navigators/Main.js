import React from 'react'
import { createDrawerNavigator } from '@react-navigation/drawer'

import ListMember from '../Containers/Members/ListMember'
import ListTask from '../Containers/Tasks/ListTask'
import ListRoles from '../Containers/Roles/ListRoles'
import ProjectInfo from '../Containers/Projects/ProjectInfo'
import Setting from '../Containers/Setting/Setting'

const Drawer = createDrawerNavigator()

const MainNavigator = () => {
  return (
    <Drawer.Navigator initialRouteName="Home">
      <Drawer.Screen name="Task" component={ListTask} />
      <Drawer.Screen name="Member" component={ListMember} />
      <Drawer.Screen name="Role" component={ListRoles} />
      <Drawer.Screen name="Project" component={ProjectInfo} />
      <Drawer.Screen name="Setting" component={Setting} />
    </Drawer.Navigator>
  )
}

export default MainNavigator
